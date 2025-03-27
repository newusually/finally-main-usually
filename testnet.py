# -*- coding: gbk -*-
import time

import torch
import torch.nn as nn
import torch.nn.functional as F
from torch.utils import data

from dataloader import Test_CsvDataset


class testNet(nn.Module):

    def __init__(self):
        super(testNet, self).__init__()
        self.fc1 = nn.Linear(4 * 9, 512)
        self.fc2 = nn.Linear(512, 512)
        self.fc3 = nn.Linear(512, 1)

    def forward(self, x):
        x = self.fc1(x)
        x = F.relu(x)
        x = self.fc2(x)
        x = F.relu(x)
        x = self.fc3(x)
        x = F.relu(x)

        return x

    '''
    这段代码的目的是将模型的输出（预测值）和实际标签（真实值）从GPU转移到CPU，并将它们乘以out_max，
    这是一个缩放因子，用于将数据从归一化的状态恢复到原始的尺度。然后，计算预测增长，这是真实值和预测值之间的差值。
    outputs = outputs.cpu() * out_max：这行代码将模型的输出（预测值）从GPU转移到CPU，并乘以out_max，
    以将预测值从归一化的状态恢复到原始的尺度。
    labels = labels.cpu() * out_max：这行代码将实际标签（真实值）从GPU转移到CPU，并乘以out_max，
    以将真实值从归一化的状态恢复到原始的尺度。
    predicted_increase = labels.numpy()[0][0] - outputs.item()：这行代码计算预测增长，
    这是真实值和预测值之间的差值。
    如果你只想返回预测价格，而不返回实际价格，你可以简单地删除与labels相关的代码，
    并将predicted_increase更改为predicted_price
    '''

    def gettest(self, files, models):
        device = torch.device("cpu")
        time.sleep(0.2)
        train_dataset = Test_CsvDataset(files, 9)
        train_loader = data.DataLoader(dataset=train_dataset, batch_size=256,
                                       shuffle=True,
                                       num_workers=0)
        net = testNet.cpu(self)

        net.load_state_dict(torch.load(models, map_location="cpu"))
        criterion = nn.MSELoss()
        net.eval()

        for inputs, labels, out_max, timets in train_loader:
            inputs = inputs.to(device)
            labels = labels.to(device)
            outputs = net(inputs)
            loss = criterion(outputs, labels)
            outputs = outputs.cpu() * out_max
            predicted_price = outputs.item()

        return predicted_price
