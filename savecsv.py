# -*- coding: gbk -*-

from databuy import Databuy

if __name__ == '__main__':
    minute = '15'

    data = Databuy()

    data.getethinfo(minute)
