import { _decorator, Component, Label, Button } from 'cc';  

const { ccclass, property } = _decorator;  

@ccclass('PopupScript')  
export default class PopupScript extends Component {  
    @property(Label)  
    messageLabel: Label = null;  

    @property(Button)  
    confirmButton: Button = null;  




    setMessage(message) {  
        this.messageLabel.string = message; // 設置顯示消息  
        
    }  

    setAction(action) {  
        if (this.confirmButton) {
            this.confirmButton.node.off('click'); // 清除舊的事件
            this.confirmButton.node.on('click', action, this); // 設置新的點擊事件
    
            console.log("Confirm button action set!"); // 確認設置完成
        } else {
            console.error("Confirm button is null or not assigned!");
        }
    

    }  
}