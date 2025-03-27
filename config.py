from yacs.config import CfgNode as CN

_C = CN()
_C.DATASETS = CN()

# TRAIN
_C.BATCH_SIZE = 216
# 数据N次循环
_C.EPOCHS = 500000000
_C.PRETRAINED_MOEDLS = ''
_C.EXP = CN()
_C.EXP.PATH = 'datas'
_C.EXP.NAME = 'training'
cfg = _C
cfg.freeze()
# if __name__ == '__main__':
# print(cfg)
