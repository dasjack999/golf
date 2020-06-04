

class WsClient {
    //
    m_socket:WebSocket;
    clientId:number=-1;
    m_heartbeat:any={
        id:'S2C_HeartBeat',

    };
    constructor(){

    }
    connect(url :string,timeOut ?:number){
        timeOut=timeOut||500;
        setTimeout(()=>{
            if(this.clientId == -1){
                this.m_socket.close();
                this.m_socket=null;
            }
        },timeOut)
        let con = this.m_socket=new WebSocket(url);
        con.onopen=this.onConnected.bind(this);
        con.onclose=this.onDisconnected.bind(this);
        con.onmessage=this.onMessage.bind(this);
    }
    send(data:any){
        data.clientId=this.clientId
        this.m_socket.send(JSON.stringify(data))
    }
    onConnected(evt:Event){

    }
    onDisconnected(evt:Event){

    }
    onMessage(evt:MessageEvent){
        let cmd=JSON.parse(evt.data);
        if(cmd.id=='connect'){
            this.clientId=cmd.clientId;
            //
            setInterval(()=>{
                this.send(this.m_heartbeat)
            },1000)
        }else if(cmd.id=='S2C_HeartBeat'){

        }
    }
}