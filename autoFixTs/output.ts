import AudioClip = cc.AudioClip;

const {ccclass, property} = cc._decorator;
export class AudioInfo {
    constructor(clip:AudioClip,id : number,isMusic : boolean,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"){
        this.IsMusic = isMusic;
        this.AudioClip = clip;
        this.Id = id;
    }
    AudioClip : AudioClip = null;
    Id : number = 0;
    IsMusic : boolean = false;
}
export class MusicInfo { 
    constructor(clip:AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"){
        this.AudioClip = clip;
    }
    AudioClip : AudioClip = null;
}
@ccclass
export default class AudioManager extends cc.Component {

    private m_CurrentPlayingAudioList : Array<AudioInfo> = [];
    private MaxAudioCount : number = 10;
    private m_bIsEnableAudio : boolean = true;
    private m_bIsEnableMusic : boolean = true;
    private m_fVolume : number = 0.8;
    private m_LoadedAudioClipMap : Map<string,AudioClip> = new Map<string, AudioClip>();
    private m_CurrentMusicInfo : MusicInfo = null;

    public Init(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"){

    }
    public SetMusicEnable(status:boolean,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        this.m_bIsEnableMusic = status;
        if(!this.m_bIsEnableMusic){
            this.PauseMusic();
        }else{
            this.ResumeMusic();
        }
    }
    public SetAudioEnable(status:boolean,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3"){
        this.m_bIsEnableAudio = status;
        if(!this.m_bIsEnableAudio){
            this.StopAllAudio();
        }
    }
    public SetVolume(volume : number,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"):void{
        this.m_fVolume = volume;
        if(this.m_fVolume < 0){
            this.m_fVolume = 0;
        }else if(this.m_fVolume > 1){
            this.m_fVolume = 1;
        }
    }
    public StopAudio(clip:AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3"):void{
        for(var i=0;i<this.m_CurrentPlayingAudioList.length;i++){
            if(this.m_CurrentPlayingAudioList[i].AudioClip == clip){
                cc.audioEngine.stop(this.m_CurrentPlayingAudioList[i].Id);
                this.m_CurrentPlayingAudioList.splice(i,1);
                return;
            }
        }
    }
    public PauseAudio(clip:AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        for(var i=0;i<this.m_CurrentPlayingAudioList.length;i++){
            if(this.m_CurrentPlayingAudioList[i].AudioClip == clip){
                cc.audioEngine.pause(this.m_CurrentPlayingAudioList[i].Id);
                return;
            }
        }
    }
    public ResumeAudio(clip:AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3"):void{
        for(var i=0;i<this.m_CurrentPlayingAudioList.length;i++){
            if(this.m_CurrentPlayingAudioList[i].AudioClip == clip){
                cc.audioEngine.resume(this.m_CurrentPlayingAudioList[i].Id);
                return;
            }
        }
    }
    private m_bLastSettingEnableMusic : boolean = false;
    public PauseAll(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"){
        this.m_bLastSettingEnableMusic = this.m_bIsEnableMusic;
        this.SetMusicEnable(false);
    }
    public ResumeAll(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"){
        if(this.m_bLastSettingEnableMusic){
            this.SetMusicEnable(true); 
        }
    }
    private StopAllAudio(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3"){
        for(var i=0;i<this.m_CurrentPlayingAudioList.length;){
            if(!this.m_CurrentPlayingAudioList[i].IsMusic){
                cc.audioEngine.stop(this.m_CurrentPlayingAudioList[i].Id);
                this.m_CurrentPlayingAudioList.splice(i,1);
            }else{
                i++;
            }
        }
    }
    public PlayAudio(clip:AudioClip,isLoop ?:boolean ,volume ?:number,isMusic ?:boolean){
        if(!isMusic){
            isMusic = false;
        }
        if(!this.m_bIsEnableAudio && !isMusic){
            return;
        }
        if(!isLoop){
            isLoop = false;
        }
        let vol = this.m_fVolume;
        if(volume){
            vol = volume;
        }
        let id = cc.audioEngine.play(clip,isLoop,vol);
        this.AddClip(id,clip,isMusic);
        this.CheckRelease();
    }
    public PlayAudioByPath(filePath : string,isLoop ?:boolean,volume ?:number){
        if(!this.m_bIsEnableAudio){
            return;
        }
        let clip = this.GetClipByPath(filePath);
        if(null != clip){
            this.PlayAudio(clip,isLoop,volume);
            return;
        }
        this.LoadClip(filePath,function () {
            clip = this.GetClipByPath(filePath);
            this.PlayAudio(clip,isLoop,volume);
        }.bind(this));
    }
    public PlayMusicByPath(filePath : string,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"){
        let clip = this.GetClipByPath(filePath);
        if(null == clip){
            this.LoadClip(filePath,function () {
                clip = this.GetClipByPath(filePath);
                this.DoPlayMusic(clip);
                return;
            }.bind(this));
        }else{
            this.DoPlayMusic(clip);
        }
    }
    public PauseMusic(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        if(null != this.m_CurrentMusicInfo){
            this.PauseAudio(this.m_CurrentMusicInfo.AudioClip);
        }
    }
    public ResumeMusic(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3"):void{
        if(null != this.m_CurrentMusicInfo){
            this.ResumeAudio(this.m_CurrentMusicInfo.AudioClip);
        }
    }
    private DoPlayMusic(audioClip : AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        if(null != this.m_CurrentMusicInfo){
            if(audioClip == this.m_CurrentMusicInfo.AudioClip){
                return;
            }
            this.StopAudio(this.m_CurrentMusicInfo.AudioClip);
            this.m_CurrentMusicInfo.AudioClip = audioClip;
        }else{
            this.m_CurrentMusicInfo = new MusicInfo(audioClip);
        }

        this.PlayAudio(this.m_CurrentMusicInfo.AudioClip,true,this.m_fVolume,true);
        if(!this.m_bIsEnableMusic){
            this.scheduleOnce(function () {
                this.PauseMusic();
            }.bind(this),0);
        }
    }
    private IsAudioPlaying(id:number,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"):boolean{
        let time = cc.audioEngine.getCurrentTime(id);
        let totalTime = cc.audioEngine.getDuration(id);
        if(totalTime <= 0){
            return false;
        }
        return time < totalTime;
    }
    private GetClipByPath(path:string,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):AudioClip{
        if(this.m_LoadedAudioClipMap.has(path)){
            return this.m_LoadedAudioClipMap.get(path);
        }
        return null;
    }
    private LoadClip(filePath:string,callback : ()=>void):void{
        cc.loader.loadRes(filePath, cc.AudioClip, function (err, clip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4") {
            if(!err){
                this.OnClipLoaded(filePath,clip);
                callback();
            }else{
                console.error("error on play clip by path " + filePath);
            }
        }.bind(this));
    }
    private OnClipLoaded(path:string,clip:AudioClip,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        this.m_LoadedAudioClipMap.set(path,clip);
    }
    private AddClip(id:number,clip : AudioClip,isMusic : boolean,randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2",randomKey_3="randomValue_3",randomKey_4="randomValue_4"):void{
        for(var i=0;i<this.m_CurrentPlayingAudioList.length;i++){
            if(this.m_CurrentPlayingAudioList[i].AudioClip == clip){
                this.m_CurrentPlayingAudioList[i].Id = id;
                let elem = this.m_CurrentPlayingAudioList[i];
                this.m_CurrentPlayingAudioList.splice(i,1);
                this.m_CurrentPlayingAudioList.push(elem);
                return;
            }
        }
        this.m_CurrentPlayingAudioList.push(new AudioInfo(clip,id,isMusic));

    }
    private CheckRelease(randomKey_0="randomValue_0",randomKey_1="randomValue_1",randomKey_2="randomValue_2"){
        if(this.m_CurrentPlayingAudioList.length <= this.MaxAudioCount){
            return;
        }
        while(this.m_CurrentPlayingAudioList.length > this.MaxAudioCount){
            let info = this.m_CurrentPlayingAudioList[0];
            if(info.IsMusic){
                continue;
            }
            this.m_CurrentPlayingAudioList.splice(0,1);
            cc.audioEngine.uncache(info.AudioClip);
        }
    }

}
