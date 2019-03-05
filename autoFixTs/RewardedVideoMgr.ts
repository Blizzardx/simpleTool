import { gcfg } from "../../gcfg";
import QWEventDispatcher from "../GMDBase/QWEventDispatcher";
import QWEvent from "../GMDBase/QWEvent";
import { SDataVideoAD } from "./SDataVideoAD";
import DailyRefresh from "../GMDBase/DailyRefresh";
import LCData from "../GMDWeakNetGameBase/LCData";
import { MD } from "../GMDBase/MD";
//激励视频管理器
let g_vidioAd
let g_retryCount: number = 0
let g_Evts: QWEventDispatcher = new QWEventDispatcher()
let g_CurrChannelId = null//当前分享口id
let g_FuncCallbacks: any = {}
let g_isLoadok = false

let g_Opening = false//视频是否处于开启状态

export default class RewardedVideoMgr {
    //分享错误事件
    public static EVT_ERROR = "EvtError"

    //广告被拉起事件
    public static EVT_ADSHOW = "ADShow"

    //错误号枚举
    public static ERR_NO = {
        NotEnded: 1,//在视频尚未观看完成的时候关闭了视频
        TooMany: 2, //观看次数太多，已经没有奖励可领取了
        UnknownError: 3 //未知错误
    };

    //激励视频系统开始工作
    public static Go() {
        if (!gcfg.ADUnitId) {
            console.error("缺少 ADUnitId 配置！")
            return
        }

        //追加需要同步的参数
        {
            var SyncStrParam = gcfg["SyncStrParam"]
            for (var channelID = 1; channelID <= 10; channelID++) {
                SyncStrParam["VAD_" + channelID] = 1 //广告观看次数
            }
        }


        //注册需要日清的数据
        {

            //清空每日看广告次数
            DailyRefresh.AddRefrishData(LCData.ParamName.ADWATCH, "0")
            //生成宝箱的次数重新加载
            DailyRefresh.AddRefrishData(LCData.ParamName.OPEN_TREASURE_BOX_TIME, "0");

            for (var channelID = 1; channelID <= 10; channelID++) {
                var cinfo = SDataVideoAD.GetChannel(channelID)
                if (!cinfo) continue
                if (cinfo[SDataVideoAD.FD.I_NumType] == 1)//限制次数，每日刷新
                {
                    DailyRefresh.AddRefrishData("VAD_" + channelID, "0")
                }
            }
        }

        //创建广告对象
        g_vidioAd = wx.createRewardedVideoAd({ adUnitId: gcfg.ADUnitId })

        //挂载广告关闭事件
        g_vidioAd.onClose((res) => {
            g_Opening = false
            var isEnded = res.isEnded//是否观看完成
            RewardedVideoMgr.OnClose(isEnded)

        })

        //激励视频加载成功
        g_vidioAd.onLoad(() => {
            g_isLoadok = true
            RewardedVideoMgr._Show()
        })

        //激励视频加载失败
        g_vidioAd.onError(err => {
            g_isLoadok = false
            RewardedVideoMgr.ReTry()
        })

    }


    //获取某个视频观看口的观看次数
    public static GetADViewingtimes(channelID: number) {
        return LCData.GetNumber("VAD_" + channelID)
    }

    //获取全部视频观看口的总观看次数
    public static GetAllADViewingtimes() {
        var cList = SDataVideoAD.ChannelList
        var sum = 0
        for(var i=0;i<cList.length;i++) sum+=RewardedVideoMgr.GetADViewingtimes(cList[i])
        return sum
    }

    //进入游戏关卡时调用
    public static EnterLevel() {
        for (var channelID = 1; channelID <= 10; channelID++) {
            var cinfo = SDataVideoAD.GetChannel(channelID)
            if (!cinfo) continue
            if (cinfo[SDataVideoAD.FD.I_NumType] == 2)//限制次数，进入关卡刷新
                LCData.Set("VAD_" + channelID, "0")
        }
    }

    //激励视频是否处于打开状态
    public static get IsOpen(): boolean { return g_Opening }


    //添加事件监听
    public static addEventListener(type: string, listener, owner = null) {
        g_Evts.addEventListener(type, listener, owner)
    }

    //弹出全屏激励广告
    //channelId 分享口ID
    public static OpenAD(channelId: number) {
        //记录当前观看的分享口id
        g_CurrChannelId = channelId

        //显示激励广告
        RewardedVideoMgr.Show()
    }

    //注册功能
    public static RegFunc(funcNum, func, owner = null) {
        g_FuncCallbacks[funcNum] = [func, owner]
    }

    ///////////////////////////////////////////////////////////////////////////////
    //内部接口
    ///////////////////////////////////////////////////////////////////////////////

    static OnClose(isEnded: boolean) {
        if (!isEnded)//没有观看完成
        {
            RewardedVideoMgr.PostErrorEvt(g_CurrChannelId, RewardedVideoMgr.ERR_NO.NotEnded)
            g_CurrChannelId = null
            return
        }

        var channelInfo = SDataVideoAD.GetChannel(g_CurrChannelId)

        //增加视频观看次数
        var watchCount = RewardedVideoMgr.GetADViewingtimes(g_CurrChannelId) + 1

        //检查是否观看太多次了
        if (watchCount > channelInfo[SDataVideoAD.FD.I_Num]) {
            RewardedVideoMgr.PostErrorEvt(g_CurrChannelId, RewardedVideoMgr.ERR_NO.TooMany)
            g_CurrChannelId = null
            return
        }

        //保存观看次数
        LCData.Set("VAD_" + g_CurrChannelId, watchCount + "")

        //给奖励
        var adFunc = channelInfo[SDataVideoAD.FD.I_AdFunc]
        if (adFunc > 0) {
            RewardedVideoMgr.DoFunc(
                adFunc,
                {
                    param: channelInfo[SDataVideoAD.FD.I_Param],//功能参数
                }
            )
        }
        g_CurrChannelId = null
    }

    //执行功能
    static DoFunc(funcNum, obj: any) {
        var funcInfo = g_FuncCallbacks[funcNum]
        if (!funcInfo) return

        try {
            if (funcInfo[1])//存在owner
                funcInfo[0].call(funcInfo[1], obj)
            else
                funcInfo[0](obj)
        } catch (msg) {
            console.error(msg)
        }
    }


    //抛出错误事件
    static PostErrorEvt(channelId: number, errNo: number, exAttrs: any = null) {
        var evt = new QWEvent(RewardedVideoMgr.EVT_ERROR)
        evt.details = {
            channelId: channelId,//分享口ID
            errNo: errNo,//错误号
        }

        //填充扩展属性
        if (exAttrs) {
            for (var k in exAttrs) evt.details[k] = exAttrs[k]
        }

        g_Evts.dispatchEvent(evt)
    }



    //显示激励广告
    static Show() {
        g_retryCount = 0
        RewardedVideoMgr._Show()
    }



    static _Show() {
        if (!g_isLoadok)//没有装载ok
        {
            g_vidioAd.load()
            return
        }

        if (!g_CurrChannelId)//没有指定分享口id
            return

        g_Opening = true

        //抛出广告被拉起事件
        var evt = new QWEvent(RewardedVideoMgr.EVT_ADSHOW)
        evt.details = {
            channelId: g_CurrChannelId,//分享口ID
        }
        g_Evts.dispatchEvent(evt)

        //显示广告
        g_vidioAd.show()
    }

    //重试
    static ReTry() {
        /*禁用重试次数，感觉小游戏底层有重试逻辑
        if(g_retryCount++<3)
        {
            console.log("重新拉取视频")  
            g_vidioAd.load()
        }
        else
        */
        {
            console.error("拉取视频失败！")
            g_Opening = false
            RewardedVideoMgr.PostErrorEvt(g_CurrChannelId, RewardedVideoMgr.ERR_NO.UnknownError)
        }
    }
}



//1号功能为加金币，直接实现在底层框架
RewardedVideoMgr.RegFunc(
    1,
    (obj: any) => {
        var currGold = LCData.GetNumber(LCData.ParamName.GOLD)
        LCData.Set(LCData.ParamName.GOLD, currGold + obj.param)
    }
)

MD.RewardedVideoMgr = RewardedVideoMgr