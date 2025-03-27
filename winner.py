# -*- coding: gbk -*-

import datetime
import multiprocessing
import time
from datetime import datetime

from apscheduler.schedulers.blocking import BlockingScheduler

from databuy import Databuy, SendDingding
from mvc import MVC


class Datainfo:
    class mywindows_multiprocessing():

        # 运行函数 必须写
        def run():
            sch = Datainfo.windowsshow()

            p1 = multiprocessing.Process(target=sch.okx3M_buy)
            p2 = multiprocessing.Process(target=sch.okx5M_buy)
            p3 = multiprocessing.Process(target=sch.okx15M_buy)
            p4 = multiprocessing.Process(target=sch.okxbuy)

            p1.start()
            p2.start()
            p3.start()
            p4.start()

            p1.join()
            p2.join()
            p3.join()
            p4.join()

    class windowsshow():

        def okx1H_buy(self):
            self.getdatainfo('1H')
            '''
            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['1H'], hour='*/1')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()
            '''

        def okx2H_buy(self):
            self.getdatainfo('2H')
            '''
            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['2H'], hour='*/2')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()
            '''

        def okx4H_buy(self):
            self.getdatainfo('4H')
            '''
            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['4H'], hour='*/4')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()
            '''

        def okx3M_buy(self):
            # self.getdatainfo('3m')

            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['3m'], minute='*/3')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()

        def okx5M_buy(self):
            # self.getdatainfo('5m')

            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['5m'], minute='*/5')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()

        def okx15M_buy(self):
            # self.getdatainfo('15m')

            scheduler = BlockingScheduler()
            scheduler.add_job(self.getdatainfo, 'cron', args=['15m'], minute='*/15')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()

        def okxbuy(self):

            # self.buysell()
            scheduler = BlockingScheduler()
            scheduler.add_job(self.buysell, 'cron', args=['15m'], minute='*/13')
            print(scheduler.get_jobs())
            try:
                scheduler.start()
            except KeyboardInterrupt:
                scheduler.shutdown()

        def buysell(self, minute):

            data = Databuy()

            data.buysell(minute)

        def getdatainfo(self, minute):

            time.sleep(1)
            print("\nminute-->>", minute, "start-->>", datetime.now())
            if (minute == '1m'):
                minute = '1'
            if (minute == '3m'):
                minute = '3'
            if (minute == '5m'):
                minute = '5'
            if (minute == '15m'):
                minute = '15'

            data = Databuy()

            symbollist = MVC.getsymbollist()

            if len(symbollist) < 10:
                SendDingding.sender("symbol null!!!--->>>" + "，\n我们是守护者，也是一群时刻对抗危险和疯狂的可怜虫！！！")

            p1 = multiprocessing.Process(target=data.getbuyinfo, args=[symbollist[:50], minute])
            p2 = multiprocessing.Process(target=data.getbuyinfo, args=[symbollist[50:100], minute])
            p3 = multiprocessing.Process(target=data.getbuyinfo, args=[symbollist[100:150], minute])
            p4 = multiprocessing.Process(target=data.getbuyinfo, args=[symbollist[150:], minute])

            p1.start()
            p2.start()
            p3.start()
            p4.start()

            p1.join()
            p2.join()
            p3.join()
            p4.join()

            print("\nminute-->>", minute, "end-->>", datetime.now())


if __name__ == '__main__':
    # 启动

    print('=============================================================================================')
    print('我们是守护者，也是一群时刻对抗危险和疯狂的可怜虫 ！^_^')
    print('=============================================================================================')
    print('=============================================================================================')

    mywindowsmultiprocessing = Datainfo.mywindows_multiprocessing
    mywindowsmultiprocessing.run()
