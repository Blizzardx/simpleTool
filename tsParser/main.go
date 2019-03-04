package main

import (
	"github.com/robertkrimen/otto/parser"
	"fmt"
)

func main() {
	filename := "" // A filename is optional
	src := `
(function() {"use strict";var __module = CC_EDITOR ? module : {exports:{}};var __filename = 'preview-scripts/assets/script/project/Sdk/SdkManager.js';var __require = CC_EDITOR ? function (request) {return cc.require(request, require);} : function (request) {return cc.require(request, __filename);};function __define (exports, require, module) {"use strict";
cc._RF.push(module, '1d30f31gx9Ic52RMfSd6z5U', 'SdkManager', __filename);
// script/project/Sdk/SdkManager.ts

Object.defineProperty(exports, "__esModule", { value: true });
var NetworkManager_1 = require("../../common/NetworkManager");
var Common_1 = require("../../common/Common");
var SdkShareInfo = /** @class */ (function () {
    function SdkShareInfo(title, image, desc) {
        this.ShareTitle = "";
        this.ShareImage = "";
        this.ShareDesc = "";
        this.ShareImage = image;
        this.ShareTitle = title;
        this.ShareDesc = desc;
    }
    return SdkShareInfo;
}());
exports.SdkShareInfo = SdkShareInfo;
var SdkSetting = /** @class */ (function () {
    function SdkSetting() {
        this.ShareInfoList = [];
        this.OpenShare = false;
        this.OpenGift = false;
    }
    return SdkSetting;
}());
exports.SdkSetting = SdkSetting;
var ShareInfo = /** @class */ (function () {
    function ShareInfo(title, image) {
        this.title = "";
        this.imgUrl = "";
        this.imgUrl = image;
        this.title = title;
    }
    return ShareInfo;
}());
exports.ShareInfo = ShareInfo;
var SdkMoreGameRedirectInfo = /** @class */ (function () {
    function SdkMoreGameRedirectInfo() {
        this.app_id = "";
        this.name = "";
        this.img_url = "";
        this.path = "";
        this.param = "";
    }
    return SdkMoreGameRedirectInfo;
}());
exports.SdkMoreGameRedirectInfo = SdkMoreGameRedirectInfo;
var SdkMoreGameDateInfo = /** @class */ (function () {
    function SdkMoreGameDateInfo() {
        this.hz_app_id = "";
        this.hz_path = "";
        this.redirect = [];
    }
    return SdkMoreGameDateInfo;
}());
exports.SdkMoreGameDateInfo = SdkMoreGameDateInfo;
var SdkMoreGameInfo = /** @class */ (function () {
    function SdkMoreGameInfo() {
        this.errno = 0;
        this.errmsg = "";
        this.data = null;
    }
    return SdkMoreGameInfo;
}());
exports.SdkMoreGameInfo = SdkMoreGameInfo;
var SdkManager = /** @class */ (function () {
    function SdkManager() {
        this.m_DefaultSetting = null;
        this.m_SdkSetting = null;
        this.m_SdkMoreGameRec = null;
        this.m_SdkMoreGameTry = null;
        this.m_bIsReport = false;
    }
    SdkManager.GetInstance = function () {
        if (null == this.g_Instance) {
            this.g_Instance = new SdkManager();
            this.g_Instance.Init();
        }
        return this.g_Instance;
    };
    SdkManager.prototype.Init = function () {
        this.m_DefaultSetting = new SdkSetting();
        this.m_DefaultSetting.ShareInfoList = [];
        this.m_DefaultSetting.ShareInfoList.push(new SdkShareInfo("史上最好玩的打枪游戏!!!", "https://gather.51weiwan.com/uploads/app/20181212/4cbeabe9f796307527e589dbff5944a5.png", "右上角转发"));
        this.m_DefaultSetting.OpenShare = false;
        this.m_DefaultSetting.OpenGift = false;
        NetworkManager_1.default.GetInstance().HttpGet("https://gather.51weiwan.com/api/app/getConfig?game_id=155", this.OnCallback.bind(this));
        this.InitGameList();
        this.InitGameList();
    };
    SdkManager.prototype.OnCallback = function (isOk, result) {
        if (!isOk) {
            this.OnError("server not response");
            return;
        }
        var errorCode = result["errno"];
        if (errorCode != 0) {
            this.OnError(result["errmsg"]);
            return;
        }
        result = result["data"];
        this.m_SdkSetting = new SdkSetting();
        this.m_SdkSetting.ShareInfoList = [];
        this.m_SdkSetting.OpenShare = result["force_share"] == 1;
        this.m_SdkSetting.OpenGift = result["new_gift"] == 1;
        console.warn("force share", this.m_SdkSetting.OpenShare);
        console.warn("new_gift", this.m_SdkSetting.OpenGift);
        var shareInfo = result["share"];
        for (var key in shareInfo) {
            var shareElemInfo = shareInfo[key]["info"];
            var shareElemInfoDesc = shareInfo[key]["description"];
            if (!shareElemInfo || !shareElemInfoDesc) {
                console.error("error on pareser json obj on get response", shareInfo[key]);
                continue;
            }
            var shareDate = new SdkShareInfo(shareElemInfo["share_title"], shareElemInfo["share_img"], shareElemInfoDesc);
            this.m_SdkSetting.ShareInfoList.push(shareDate);
        }
        if (this.m_SdkSetting.ShareInfoList.length == 0) {
            this.m_SdkSetting = null;
            this.OnError("share list is empty");
            return;
        }
        console.log("get setting info succeed ", this.m_SdkSetting);
        // GameManager.GetInstance<GameManager>().GetEventManager().Dispatch(ProjectDefine.EventDefine.OnFetchGameSetting);
    };
    SdkManager.prototype.RandomGetShareInfo = function () {
        var shareInfoList = this.GetShareInfoList();
        var randomIndex = Common_1.default.getRandomRange(0, shareInfoList.length);
        var info = shareInfoList[randomIndex];
        return new ShareInfo(info.ShareTitle, info.ShareImage);
    };
    SdkManager.prototype.OnError = function (errMsg) {
        //
        console.error("error on get sdk setting ", errMsg);
        console.error("begin use default setting");
        return;
    };
    SdkManager.prototype.CreateRedirectInfo = function (name, appId, path, icon) {
        var res = new SdkMoreGameRedirectInfo();
        res.path = path;
        res.app_id = appId;
        res.name = name;
        res.img_url = icon;
        return res;
    };
    SdkManager.prototype.InitGameList = function () {
        this.m_SdkMoreGameRec = new SdkMoreGameInfo();
        this.m_SdkMoreGameRec.data = new SdkMoreGameDateInfo();
        this.m_SdkMoreGameRec.data.redirect = [];
        this.m_SdkMoreGameRec.data.redirect.push(this.CreateRedirectInfo("糖心炮弹", "wxd16915542d3df546", "/pages/index/index?gdt_vid=kuaiyan&weixinadinfo=0001&channel=kuaiyan.h5txpd.3308", "http://gather.51weiwan.com/uploads/file/20181129/a3277708ac1ae06a7c75b32abbe24643.png"));
        this.m_SdkMoreGameRec.data.redirect.push(this.CreateRedirectInfo("热血修仙", "wx20194e7827347870", "?chid=0001&subchid=kdtf", "http://gather.51weiwan.com/uploads/file/20181127/851352ee4b4be1501d1a797e3bd4b2aa.png"));
        this.m_SdkMoreGameRec.data.redirect.push(this.CreateRedirectInfo("无双三国志", "wx4ebc6d7cac323719", "?tid=JDqVtO_72AG5U", "http://gather.51weiwan.com/uploads/file/20181127/87e37fe27a4c1a73a88e7698de2b9295.png"));
        this.m_SdkMoreGameRec.data.redirect.push(this.CreateRedirectInfo("疾速飞车", "wx3c78da89feb706f2", "index.html?channel=xqyx57", "http://gather.51weiwan.com/uploads/file/20180928/97fc933d7df392cdc07a2ec1cddb2b29.jpg"));
        console.log("On Get More Game Info ", this.m_SdkMoreGameRec);
        this.m_SdkMoreGameTry = new SdkMoreGameInfo();
        this.m_SdkMoreGameTry.data = new SdkMoreGameDateInfo();
        this.m_SdkMoreGameTry.data.redirect = [];
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("绝地碰碰车", "wxd844e5868051fbd1", "?terrace_id=37&td_id=53", "http://gather.51weiwan.com/uploads/file/20180910/63b46ee6e7af95e29009006007bbabb9.png"));
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("六角消消乐", "wx2668677937e0fda6", "?terrace_id=37&td_id=62", "http://gather.51weiwan.com/uploads/file/20180823/bce7e814289172f1162252c037dfc156.jpg"));
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("水果飞镖", "wx5e802405ef11635b", "?terrace_id=37&td_id=60", "http://gather.51weiwan.com/uploads/file/20180928/65505b3da2457ba5799fd738caaa4d67.jpg"));
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("玩命旅途", "wx8c2c700e6794bbe9", "?terrace_id=37&td_id=61", "http://gather.51weiwan.com/uploads/file/20180928/b60346473645b503e48bbcb2ed54c9e2.jpg"));
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("巨网游戏", "wx9e9b545ce617a380", "", "http://gather.51weiwan.com/uploads/file/20180822/0ba0d9174433ac652746197c36111230.png"));
        this.m_SdkMoreGameTry.data.redirect.push(this.CreateRedirectInfo("双拼消消乐", "wx646c223661383d08", "?terrace_id=37&td_id=68", "http://gather.51weiwan.com/uploads/file/20180823/0aa1b011279c8d3bf288e93255a59892.jpg"));
        // GameManager.GetInstance<GameManager>().GetEventManager().Dispatch(ProjectDefine.EventDefine.OnFetchMoreRecGame);
        console.log("On Get More Game Info ", this.m_SdkMoreGameTry);
    };
    SdkManager.prototype.GetShareInfoList = function () {
        if (this.m_SdkSetting != null) {
            return this.m_SdkSetting.ShareInfoList;
        }
        return this.m_DefaultSetting.ShareInfoList;
    };
    SdkManager.prototype.GetMoreGameTryListRandom = function () {
        if (null == this.m_SdkMoreGameRec) {
            return null;
        }
        if (!this.m_SdkMoreGameRec.data || !this.m_SdkMoreGameRec.data.redirect) {
            return null;
        }
        var result = [];
        var randomRes = [];
        for (var i = 0; i < this.m_SdkMoreGameRec.data.redirect.length; i++) {
            randomRes.push(i);
        }
        for (var i = 0; i < this.m_SdkMoreGameRec.data.redirect.length; i++) {
            var randomValue = Common_1.default.getRandomRange(0, randomRes.length);
            var randomIndex = randomRes[randomValue];
            randomRes.splice(randomValue, 1);
            result.push(this.m_SdkMoreGameRec.data.redirect[randomIndex]);
        }
        return result;
    };
    //获取推荐游戏列表
    SdkManager.prototype.GetMoreGameTryList = function () {
        if (null == this.m_SdkMoreGameRec) {
            return null;
        }
        if (!this.m_SdkMoreGameRec.data || !this.m_SdkMoreGameRec.data.redirect) {
            return null;
        }
        return this.m_SdkMoreGameRec.data.redirect;
    };
    //获取可跳转游戏列表
    SdkManager.prototype.GetMorrGameRecList = function () {
        if (null == this.m_SdkMoreGameTry) {
            return null;
        }
        if (!this.m_SdkMoreGameTry.data || !this.m_SdkMoreGameTry.data.redirect) {
            return null;
        }
        return this.m_SdkMoreGameTry.data.redirect;
    };
    SdkManager.prototype.OpenGame = function (appId, path) {
        if (cc.sys.platform != cc.sys.WECHAT_GAME) {
            return;
        }
        window["wx"].navigateToMiniProgram({
            appId: appId,
            path: path,
            envVersion: 'release',
            success: function (res) {
                console.log("打开成功：" + res);
            },
            fail: function (res) {
                console.log("打开失败：" + res);
            }
        });
    };
    SdkManager.prototype.OnClickMoreGame = function () {
        this.OpenGame("wx9e9b545ce617a380", "");
    };
    SdkManager.prototype.GetRandomGameInfo = function () {
        var gameList = this.GetMorrGameRecList();
        if (!gameList || gameList == null || gameList.length == 0) {
            return null;
        }
        var randomIndex = Common_1.default.getRandomRange(0, gameList.length);
        return gameList[randomIndex];
    };
    SdkManager.prototype.GetIsForceShare = function () {
        if (this.m_SdkSetting) {
            return this.m_SdkSetting.OpenShare;
        }
        if (this.m_DefaultSetting) {
            return this.m_DefaultSetting.OpenShare;
        }
        return false;
    };
    SdkManager.prototype.GetIsViewAdEnable = function () {
        return false;
    };
    SdkManager.prototype.GetIsEnableShareOrViewAd = function () {
        return this.GetIsForceShare() || this.GetIsViewAdEnable();
    };
    SdkManager.g_Instance = null;
    return SdkManager;
}());
exports.default = SdkManager;

cc._RF.pop();
        }
        if (CC_EDITOR) {
            __define(__module.exports, __require, __module);
        }
        else {
            cc.registerModuleFunc(__filename, function () {
                __define(__module.exports, __require, __module);
            });
        }
        })();
        //# sourceMappingURL=SdkManager.js.map
        
`

	// Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
	program, err := parser.ParseFile(nil, filename, src, 0)
	fmt.Printf("%v",program,err)
}
