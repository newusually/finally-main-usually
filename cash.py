# -*- coding: gbk -*-
from mvc import MVC
from userinfo import User

if __name__ == '__main__':
    api_key, secret_key, passphrase, flag = User.get_userinfo()

    #��ȡʵʱ�˻��ʽ���Ϣ ÿ���Ӳ�ѯһ��
    MVC.getcashbal(api_key, secret_key, passphrase, flag)

