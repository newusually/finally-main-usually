# -*- coding: gbk -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()


    ##����ʷ��¼ ÿ���Ӽ�¼һ��
    MVC.getcashhistory(api_key, secret_key, passphrase, flag)
