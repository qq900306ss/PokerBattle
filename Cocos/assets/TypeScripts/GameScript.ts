import { _decorator, Component, Node ,sys ,Label ,EditBox, Button,director ,Sprite , resources ,SpriteFrame }  from 'cc';
const { ccclass, property } = _decorator;

@ccclass('GameScript')
export class GameScript extends Component {


    @property(Label)  
    Connectext: Label = null;  

    @property(Label)  
    OpponentName: Label = null;  

    @property(Label)  
    UserInfo: Label = null; 

    @property(Sprite)  
    MyCard: Sprite = null; 

    @property(Sprite)  
    OpponentCard : Sprite = null; 

    @property(EditBox)  
    BetEdit: EditBox = null; 

    @property(Button)  
    BetButton: Button = null;  


    @property(Button)  
    GiveUpButton: Button = null;  

    private websocket: WebSocket = null;  



     cardMapping = {
        '黑桃': 'A',
        '紅心': 'B',
        '方塊': 'C',
        '梅花': 'D'
    };


    start() {

        // 建立 WebSocket 連接  
        this.websocket = new WebSocket('ws://localhost:8080/ws');  
        console.log('開始WebSocket connection established');  
        this.OpponentName.node.active = false;  
        this.MyCard.node.active = false;  
        this.OpponentCard.node.active = false;  
        this.BetEdit.node.active = false;  
        this.BetButton.node.active = false;  
        this.GiveUpButton.node.active = false;  


        let username = sys.localStorage.getItem('username');  
        let money = parseInt(   sys.localStorage.getItem('money') );  
        
        // 當連接打開時  
        this.websocket.onopen = () => {  
            console.log('WebSocket connection established');  

            // 取得用戶名並發送到後端  
            

            const message = JSON.stringify({  Name: username , Money : (money) });  
            this.websocket.send(message);  
            this.UserInfo.string = "此玩家帳號:" + username + " \n金幣:" + money;  
            console.log('Sent message:', message);  
        };  



        
        // 當接收到消息  
        this.websocket.onmessage = (event) => {  
            if (JSON.parse(event.data).Event === "對手ID"){
                this.OpponentName.string = "對手帳號:" + JSON.parse(event.data).Name;  
                this.OpponentName.node.active = true;  

                this.Connectext.string = "即將開始:" ; 

            }


            if (JSON.parse(event.data).Event === "可以配卡了"){
                this.Connectext.node.active = false;  
                this.MyCard.node.active = true;  
                console.log("可以配卡了:" + JSON.parse(event.data).Card.suit + "沒.suit" + JSON.parse(event.data).Card);  
                
                let newImagePath =   this.cardMapping[JSON.parse(event.data).Card.Suit] + JSON.parse(event.data).Card.Value + "/spriteFrame"



                 // 獲取 MyCard 的 Sprite 組件  
                    // 載入新的 SpriteFrame  
                    resources.load(newImagePath,SpriteFrame, (err, spriteFrame) => {  
                        if (err) {  
                            console.error(err);  
                            return;  
                        }  
                        // 將新的 SpriteFrame 設置給 Sprite 組件  
                        this.MyCard.spriteFrame = spriteFrame;  
                    });
                        


                this.BetEdit.node.active = true;  
                this.BetButton.node.active = true;
                this.GiveUpButton.node.active = true;  

                this.BetButton.node.on('click', () => {  
                    let inputValue = parseInt(this.BetEdit.string);  
                    if (isNaN(inputValue) || inputValue <= 10){

                        console.log("金額沒大於10");  
                        this.BetEdit.string = "";  
                    }else{
                        console.log("下注金額:" + inputValue);  
                        const message = JSON.stringify({ Name: username , Money : inputValue})

                        this.BetEdit.node.active = false;  
                        this.BetButton.node.active = false;  
                        this.Connectext.string = "等待對手下注中....";  
                        this.Connectext.node.active = true;  


                        this.websocket.send(message)
                        console.log('Sent 下注message:', message);  
                    }
                })
                
                this.GiveUpButton.node.on('click', () => {  
                    const message = JSON.stringify({ Name: "" , Money : 0})
                    this.websocket.send(message)
                    money = money - 10;

                    sys.localStorage.setItem('money', money);

                    this.Connectext.string = "你放棄了但還是等待對手下注中....";  
                    this.BetEdit.node.active = false;  
                    this.Connectext.node.active = true;  


                })
            }


            if (JSON.parse(event.data).Event === "你放棄下注"){

                
                director.loadScene("A"); // 切換到場景 A  


                }

            if (JSON.parse(event.data).Event === "對手放棄下注了"){

                this.Connectext.string = "對手放棄下注了,三秒後回到選單"; 
                this.GiveUpButton.node.active = false;  

                director.pause();
                setTimeout(function() {
                    // 3秒後恢復遊戲
                    director.resume();
                    director.loadScene("A"); // 切換到場景 A  

                }, 3000);


            }

            


            if (JSON.parse(event.data).Event === "收到對手傳送卡片資料"){
                this.OpponentCard.node.active = true;  
                this.GiveUpButton.node.active = false;  

                console.log("對手卡牌:" + JSON.parse(event.data).Card.Suit + "沒.suit" + JSON.parse(event.data).Card);  


                let newImagePath =   this.cardMapping[JSON.parse(event.data).Card.Suit] + JSON.parse(event.data).Card.Value + "/spriteFrame"

                 // 獲取 MyCard 的 Sprite 組件  
                    // 載入新的 SpriteFrame  
                    resources.load(newImagePath,SpriteFrame, (err, spriteFrame) => {  
                        if (err) {  
                            console.error(err);  
                            return;  
                        }  
                        // 將新的 SpriteFrame 設置給 Sprite 組件  
                        this.OpponentCard.spriteFrame = spriteFrame;  
                    });

            }

            if (JSON.parse(event.data).Event === "你贏了"){
                
            
                money = money + parseInt(this.BetEdit.string) * 2;
                sys.localStorage.setItem('money', money);

                this.Connectext.string = JSON.parse(event.data).Event + "5秒後切換回選單";  
                this.Connectext.node.active = true;  


                director.pause();
                setTimeout(function() {
                    // 2秒後恢復遊戲
                    director.resume();
                    director.loadScene("A"); // 切換到場景 A  

                }, 5000);

            }

            if (JSON.parse(event.data).Event === "你輸了"){
                
                
                money = money - parseInt(this.BetEdit.string)
                sys.localStorage.setItem('money', money);

                this.Connectext.string = JSON.parse(event.data).Event  + "5秒後切換回選單";  
                this.Connectext.node.active = true;  

                director.pause();
                setTimeout(function() {
                    // 2秒後恢復遊戲
                    director.resume();
                    director.loadScene("A"); // 切換到場景 A  

                }, 5000);

            }


            if (JSON.parse(event.data).Event === "平手"){
                

                this.Connectext.string = JSON.parse(event.data).Event + "5秒後切換回選單";  
                this.Connectext.node.active = true;  

                director.pause();
                setTimeout(function() {
                    // 2秒後恢復遊戲
                    director.resume();
                    director.loadScene("A"); // 切換到場景 A  

                }, 5000);
            }


        }
                
        // 當連接關閉時  
        this.websocket.onclose = () => {  
            console.log('WebSocket connection closed');  
        };  

        // 當發生錯誤時  
        this.websocket.onerror = (error) => {  
            console.error('WebSocket error:', error);  
        };  

    }

}


