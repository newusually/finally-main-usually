# -*- coding: gbk -*-
import csv, torch
import torch.utils.data as data
import numpy as np
import pandas as pd
import multiprocessing
import warnings

warnings.filterwarnings('ignore')
pd.set_option('mode.chained_assignment', None)


class CsvDataset(data.Dataset):

    def __init__(self, path, step):
        inputs = []
        outputs = []
        with open(path, 'r', encoding='UTF-8') as (f):
            f_c = csv.reader(f)
            idx = 0
            for row in f_c:
                if idx == 0:
                    idx += 1
                    continue
                tmp = row[1:]
                # print(tmp)
                tmp = [float(i) for i in tmp]
                outputs.append(tmp.pop(0))
                inputs.append(tmp)
                idx += 1

        self.inputs = np.array(inputs)
        self.outputs = np.array(outputs)[:, np.newaxis]
        self.step = step
        self.in_max = np.abs(self.inputs).max(0)
        self.out_max = np.abs(self.outputs).max(0)
        # print(self.in_max)
        # print(self.out_max)

    def __len__(self):
        return self.inputs.shape[0] - self.step - 1

    def __getitem__(self, index):
        _input = self.inputs[index:index + self.step]
        output = self.outputs[index + self.step + 1:index + self.step + 2]
        _input = _input / self.in_max
        output = output / self.out_max
        _input = torch.Tensor(_input)
        output = torch.Tensor(output)
        _input = _input.reshape(-1)
        output = output.reshape(-1)
        return (_input, output)


class Test_CsvDataset(data.Dataset):

    def __init__(self, path, step):
        inputs = []
        outputs = []
        time = []
        with open(path, 'r', encoding='gbk') as (f):
            f_c = csv.reader(f)
            idx = 0
            for row in f_c:
                if idx == 0:
                    idx += 1
                    continue
                tmp = row[1:]
                tmp = [float(i) for i in tmp]
                tmp = [float(i) for i in tmp]
                time.append(row[:1])
                outputs.append(tmp.pop(0))
                inputs.append(tmp)
                idx += 1

        self.inputs = np.array(inputs)
        self.outputs = np.array(outputs)[:, np.newaxis]
        self.time = time
        self.step = step
        self.in_max = np.abs(self.inputs).max(0)
        self.out_max = np.abs(self.outputs).max(0)
        # print(self.outputs.shape)
        self.inputs = self.inputs[-self.step:]
        self.outputs = self.outputs[-1:]
        self.time = self.time[-1:]

    def __len__(self):
        return 1

    def __getitem__(self, index):
        # print(np.transpose(self.inputs[0][0]))
        # print(np.transpose(self.outputs))
        _input = self.inputs[index:index + self.step]
        output = self.outputs
        _input = _input / self.in_max
        output = output / self.out_max
        _input = torch.Tensor(_input)
        output = torch.Tensor(output)
        _input = _input.reshape(-1)
        output = output.reshape(-1)
        # print((_input.cpu()[0] * self.in_max)[:9])
        # print((output.cpu() * self.out_max)[:9])
        return (_input, output, self.out_max, self.time)
