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
    ��δ����Ŀ���ǽ�ģ�͵������Ԥ��ֵ����ʵ�ʱ�ǩ����ʵֵ����GPUת�Ƶ�CPU���������ǳ���out_max��
    ����һ���������ӣ����ڽ����ݴӹ�һ����״̬�ָ���ԭʼ�ĳ߶ȡ�Ȼ�󣬼���Ԥ��������������ʵֵ��Ԥ��ֵ֮��Ĳ�ֵ��
    outputs = outputs.cpu() * out_max�����д��뽫ģ�͵������Ԥ��ֵ����GPUת�Ƶ�CPU��������out_max��
    �Խ�Ԥ��ֵ�ӹ�һ����״̬�ָ���ԭʼ�ĳ߶ȡ�
    labels = labels.cpu() * out_max�����д��뽫ʵ�ʱ�ǩ����ʵֵ����GPUת�Ƶ�CPU��������out_max��
    �Խ���ʵֵ�ӹ�һ����״̬�ָ���ԭʼ�ĳ߶ȡ�
    predicted_increase = labels.numpy()[0][0] - outputs.item()�����д������Ԥ��������
    ������ʵֵ��Ԥ��ֵ֮��Ĳ�ֵ��
    �����ֻ�뷵��Ԥ��۸񣬶�������ʵ�ʼ۸�����Լ򵥵�ɾ����labels��صĴ��룬
    ����predicted_increase����Ϊpredicted_price
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
