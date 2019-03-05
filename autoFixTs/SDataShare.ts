import * as Share from "../../sdata/res/Share"
import {ShareSimpleData} from "../../sdata/res/ShareSimple"

import jscsv from "../GMDBase/jscsv"
import Environment from "../GMDBase/Environment";

class ShareChannelInfo
{
    public NumType:number//分享口类型
    public ID:number//分享口ID
    public Rows:any//本分享口内的行
    public MaxNum:number//最大的编号

    constructor()
    {
        this.MaxNum = 0
        this.Rows = {}
    }

    //存档位
    public get SaveChannel():number
    {
        for(var key in this.Rows)
        { 
            var row = this.Rows[key]
            if(SDataShare.Singletion.FD.I_SaveChannel)//定义了存档位
                return  row[SDataShare.Singletion.FD.I_SaveChannel]
            else//没有定义存档位，采用分享口作为存档位
                return  row[SDataShare.Singletion.FD.I_ShareChannel]
        }
        
    }


    //根据今日分享次数获取一行
    public GetRow(shareNum:number)
    {
        shareNum++
        if(this.NumType==1)//根据分享次数
        {
            var row = this.Rows[shareNum]
            if(!row) row = this.Rows[0]
            if(!row) row = this.Rows[this.MaxNum]
            return row
        }

 

        //根据权重
        //计算总权重
        var tqz = 0
        for(var k in this.Rows)
        {
            var row = this.Rows[k]
            var num = row[SDataShare.Singletion.FD.I_Number]

            if(num<=0) continue
            tqz+=num
        }


        var rqz = Math.floor(Math.random()*tqz)
        var tmpqz = 0
        for(var k in this.Rows)
        {
            var row = this.Rows[k]
            var num = row[SDataShare.Singletion.FD.I_Number]
            if(num<=0) continue

            tmpqz+=num
            if(tmpqz>rqz) return row
        }
        console.error("ShareChannelInfo.GetRow 权重随机未选中任何行,返回0")
        return this.Rows[0]
        
    }
}

export default class SDataShare extends jscsv
{
    static m_Singletion:SDataShare = null
    public static get Singletion() {
        if(!SDataShare.m_Singletion ) SDataShare.m_Singletion = new SDataShare()
        return SDataShare.m_Singletion 
    }


    constructor()
    {
        super(Environment.SupportShare?Share.ShareData:ShareSimpleData)
        
        var ShareChannel = {}

        this.Foreach((key,row)=>{
            var channel = row[this.FD.I_ShareChannel]
            var num = row[this.FD.I_Number]
            if(ShareChannel[channel]) 
            { 
                var obj = ShareChannel[channel]
                obj.Rows[  num ] = row
                if(num>obj.MaxNum) obj.MaxNum = num
                return 
            }

            var attr = new ShareChannelInfo()
            attr.NumType = row[this.FD.I_NumType]//分享口类型
            attr.ID = channel//分享口ID 
            attr.Rows[ num ] = row
            if(num>attr.MaxNum) attr.MaxNum = num

            ShareChannel[channel] = attr
        })

        this.m_ShareChannel = ShareChannel 
    }

    //根据分享口和分享次数获取一行
    public GetShareRow(channel:number,shareNum:number)
    {
        if(!this.m_ShareChannel){
            return 0;
        }
        return this.m_ShareChannel[channel].GetRow(shareNum)
    }

    //获取某个分享口 每天最大分享次数
    public GetShareTotalNum(channel:number)
    {
        if(!this.m_ShareChannel){
            return 20;
        }
        return this.m_ShareChannel[channel].MaxNum
    }

    
    //获取某个分享口 存档位
    public GetSaveChannel(channel:number)
    {
        if(!this.m_ShareChannel){
            return 0;
        }
        return this.m_ShareChannel[channel].SaveChannel
    }
    

    m_ShareChannel:any
}
 