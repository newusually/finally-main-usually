import datetime
import numpy as np
import talib
import pandas as pd
import torch
import torch.nn as nn
from torch.utils.data import Dataset, DataLoader
from sklearn.preprocessing import MinMaxScaler
from sklearn.model_selection import train_test_split

class RNNModel(nn.Module):
    def __init__(self, input_size, hidden_size, num_layers, output_size):
        super(RNNModel, self).__init__()
        self.hidden_size = hidden_size
        self.num_layers = num_layers
        self.rnn = nn.RNN(input_size, hidden_size, num_layers, batch_first=True)
        self.fc = nn.Linear(hidden_size, output_size)
    def forward(self, x):
        h0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size).cuda()
        out, _ = self.rnn(x, h0)
        out = self.fc(out[:, -1, :])
        return out

class EthDataset(Dataset):
    def __init__(self, data, target):
        self.data = data
        self.target = target

    def __len__(self):
        return len(self.data)

    def __getitem__(self, index):
        return self.data[index], self.target[index]

def parse_txt_content(content):
    lines = content.strip().split('\n')
    data = {
        'eth_close': [],
        'a': [],
        'b': [],
        'e': [],
        'a-b': [],
        'b-e': [],
        'macd': [],
        'signal': [],
        'hist': [],
    }

    for line in lines:
        items = line.split(',')

        close = float(items[7].split('>>')[1])
        data['eth_close'].append(close)
        data['a'].append(float(items[2].split('>>')[1]))
        data['b'].append(float(items[3].split('>>')[1]))
        data['e'].append(float(items[4].split('>>')[1]))
        data['a-b'].append(float(items[5].split('>>')[1]))
        data['b-e'].append(float(items[6].split('>>')[1]))

    macd, signal, hist = talib.MACD(np.array(data['eth_close']), fastperiod=12, slowperiod=26, signalperiod=9)

    for i in range(len(lines)):
        data['macd'].append(macd[i])
        data['signal'].append(signal[i])
        data['hist'].append(hist[i])

    df=pd.DataFrame(data)
    df=df.dropna()
    df.to_csv("ethdata.csv",index=0)
    return df

def load_data(file_path):
    with open(file_path, 'r') as file:
        content = file.read()
    return parse_txt_content(content)

def create_dataset(df, lookback):
    data = df.values
    X, y = [], []
    for i in range(len(data) - lookback):
        X.append(data[i:i + lookback])
        y.append(data[i + lookback, 0])
    return np.array(X), np.array(y)

def train_rnn_model(df):
    lookback = 20
    input_size = df.shape[1]
    hidden_size = 64
    num_layers = 20  # 将模型层数上升到20层
    output_size = 1
    batch_size = 64
    epochs = 50

    # Split data into training and testing sets
    train_data, test_data = train_test_split(df, test_size=0.2, shuffle=False)
    train_X, train_y = create_dataset(train_data, lookback)
    test_X, test_y = create_dataset(test_data, lookback)

    # Normalize the data
    scaler_X = MinMaxScaler()
    train_X = scaler_X.fit_transform(train_X.reshape(-1, input_size)).reshape(-1, lookback, input_size)
    test_X = scaler_X.transform(test_X.reshape(-1, input_size)).reshape(-1, lookback, input_size)

    scaler_y = MinMaxScaler()
    train_y = scaler_y.fit_transform(train_y.reshape(-1, 1))
    test_y = scaler_y.transform(test_y.reshape(-1, 1))

    print(np.any(np.isnan(train_X)))
    print(np.any(np.isnan(train_y)))
    print(np.any(np.isinf(train_X)))
    print(np.any(np.isinf(train_y)))

    # Create PyTorch datasets and data loaders
    train_dataset = EthDataset(torch.tensor(train_X, dtype=torch.float32), torch.tensor(train_y, dtype=torch.float32))
    train_loader = DataLoader(train_dataset, batch_size=batch_size, shuffle=True)
    test_dataset = EthDataset(torch.tensor(test_X, dtype=torch.float32), torch.tensor(test_y, dtype=torch.float32))
    test_loader = DataLoader(test_dataset, batch_size=batch_size, shuffle=False)

    # Initialize the model, loss function, and optimizer
    model = RNNModel(input_size, hidden_size, num_layers, output_size).cuda()
    criterion = nn.MSELoss()

    optimizer = torch.optim.Adam(model.parameters(), lr=0.01)

    # Train the model
    for epoch in range(epochs):
        model.train()
        for batch_X, batch_y in train_loader:
            batch_X, batch_y = batch_X.cuda(), batch_y.cuda()
            optimizer.zero_grad()
            outputs = model(batch_X)
            loss = criterion(outputs, batch_y)
            loss.backward()
            #torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1)
            optimizer.step()

        # Evaluate the model on the test set
        model.eval()
        with torch.no_grad():
            test_loss = 0
            for batch_X, batch_y in test_loader:
                batch_X, batch_y = batch_X.cuda(), batch_y.cuda()
                outputs = model(batch_X)
                loss = criterion(outputs, batch_y)
                test_loss += loss.item()
            print(f'Epoch [{epoch + 1}/{epochs}], Test Loss: {test_loss / len(test_loader):.6f}')

    # Save the model
    torch.save(model.state_dict(), 'rnn_model.pth')

def load_rnn_model(df):
    input_size = df.shape[1]
    hidden_size = 64
    num_layers = 20  # 将模型层数上升到20层
    output_size = 1

    model = RNNModel(input_size, hidden_size, num_layers, output_size).cuda()
    model.load_state_dict(torch.load('rnn_model.pth'))
    model.eval()

    return model

if __name__ == '__main__':
    file_path = '15min_buylog.txt'
    df = load_data(file_path)

    # Train and save the RNN model
    train_rnn_model(df)

    # Load the RNN model
    model = load_rnn_model(df)

    # Get the features for correlation calculation
    features = ['eth_close', 'a', 'b', 'e', 'a-b', 'b-e', 'macd']

    # Calculate the correlations between the features and the predicted close prices
    for feature in features:
        feature_values = df[feature].values.reshape(-1, 1)

        # Use the model to make predictions
        lookback = 20
        input_size = df.shape[1]
        feature_values = np.hstack([feature_values[-lookback:], np.zeros((lookback, input_size - 1))])

        feature_values = torch.tensor(feature_values, dtype=torch.float32).unsqueeze(0).cuda()

        predictions = model(feature_values).cpu().detach().numpy()
        print(predictions)
        # Subtract the feature values from the predictions
        diff = predictions.flatten() - feature_values[:, 0, 0].cpu().numpy()
        print(diff)
        # Make sure there are no NaN values in the input arrays
        valid_indices = ~np.isnan(diff) & ~np.isnan(feature_values[:, 0, 0].cpu().numpy())

        # Calculate the correlation coefficient between the differences and the feature values
        if np.sum(valid_indices) > 1:
            correlation = np.corrcoef(diff[valid_indices], feature_values[:, 0, 0].cpu().numpy()[valid_indices])[0, 1]
        else:
            correlation = np.nan
        print(f'Correlation between predicted close prices and {feature}: {correlation:.6f}')


    import matplotlib.pyplot as plt

    # 绘制以太坊的实际收盘价和预测收盘价
    plt.figure(figsize=(12, 6))
    plt.plot(df['eth_close'], label='Actual Close Price')
    plt.plot(predictions, label='Predicted Close Price')
    plt.xlabel('Time')
    plt.ylabel('Price')
    plt.title('Actual vs Predicted Close Prices')
    plt.legend()
    plt.show()

    # 绘制特征和预测收盘价之间的关系
    for feature in features:
        plt.figure(figsize=(12, 6))
    plt.plot(df[feature], label=f'{feature} Values')
    plt.plot(predictions, label='Predicted Close Price')
    plt.xlabel('Time')
    plt.ylabel('Values')
    plt.title(f'Predicted Close Prices vs {feature}')
    plt.legend()
    plt.show()