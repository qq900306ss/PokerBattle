import { _decorator, Component, Label, Node, sys, Button, director } from 'cc';  
const { ccclass, property } = _decorator;  

@ccclass('ScaneAts')  
export class ScaneAts extends Component {  

    @property(Label)  
    userInfo: Label = null;  

    @property(Button)  
    StartButton: Button = null;  


    onLoad() {  
        this.StartButton.node.on('click', this.onStartClicked, this);  
    }  

    start() {  
        let username = sys.localStorage.getItem('username');  
        let money = sys.localStorage.getItem('money');  

        console.log("Retrieved username: ", username);  
        console.log("Retrieved money: ", money);  

        this.userInfo.string = "遊戲帳號: " + username + "  金幣: " + money;  
    }  

    onStartClicked() {  
        console.log('按下按鈕StartButton');  
        
        // 切換到場景 B  
        director.loadScene("B");  
    }  
}