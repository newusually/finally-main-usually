# -*- coding: utf-8 -*-
import sys
import datetime
import numpy as np
import talib
from PyQt5.QtWidgets import QApplication, QMainWindow, QVBoxLayout, QWidget, QLabel, QMenu, QMenuBar, QAction, QMessageBox
from PyQt5.QtCore import QTimer, QEventLoop, Qt
from matplotlib.backends.backend_qt5agg import FigureCanvasQTAgg as FigureCanvas
from matplotlib.backends.backend_qt5 import NavigationToolbar2QT as NavigationToolbar
from matplotlib.figure import Figure
import matplotlib.dates as mdates
import mplcursors
import pandas as pd

def parse_txt_content(content):
    lines = content.strip().split('\n')
    data = {
        'time': [],
        'eth_close': [],
        'a': [],
        'b': [],
        'e': [],
        'a-b': [],
        'b-e': [],
        'macd': [],
        'signal': [],
        'hist': [],
        'buy_label': [],
        'next_20_closes': [],

    }

    for line in lines:
        items = line.split(',')
        time_str = items[1].split('>>')[1]
        time_obj = datetime.datetime.strptime(time_str, '%Y-%m-%d %H:%M:%S')
        data['time'].append(time_obj)
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

        if i > 0 and 0.1<data['a'][i] < 0.2  and 0.1<data['b'][i] < 0.2  and data['e'][i] < 0.1 and 0<data['a-b'][i] < 0.1  and 0<data['b-e'][i] < 0.1 \
                and 1<data['macd'][i]< 2 :
            data['buy_label'].append(1)
        else:
            data['buy_label'].append(0)

        if i + 20 < len(lines):
            next_20_closes = ','.join([str(data['eth_close'][j]) for j in range(i+1, i+20)])
        else:
            next_20_closes = ''
        data['next_20_closes'].append(next_20_closes)


    df=pd.DataFrame(data)
    df.to_csv("ethdata.csv",index=0)


    return data



if __name__ == '__main__':
    file_path='15min_buylog.txt'
    with open(file_path, 'r') as file:
        content = file.read()
    data = parse_txt_content(content)
    print(data)