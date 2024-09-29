import { _decorator, Component, Node, EditBox,director, Button , instantiate, sys  } from 'cc';  

const { ccclass, property } = _decorator;  

@ccclass('LoginScript')  
export default class LoginScript extends Component {  

    @property(EditBox)  
    usernameInput: EditBox = null;  

    @property(EditBox)  
    passwordInput: EditBox = null;  

    @property(Button)  
    loginButton: Button = null;  

    onLoad() {  
        console.log(this.loginButton); // 檢查 loginButton 是否為 null  
        if (this.loginButton) {  
            this.loginButton.node.on('click', this.onLoginClicked, this);  
        } else {  
            console.error('loginButton is not assigned!');  
        }  
       
    }  

    onLoginClicked() {  
        let username = this.usernameInput.string;  
        let password = this.passwordInput.string;  

        if (!username || !password) {  
            console.log('Please enter both username and password');  
            return;  
        }  

        let loginData = {  
            username: username,  
            password: password  
        };  

        // 發送 HTTP POST 請求到後端 Golang  
        let xhr = new XMLHttpRequest();  
        xhr.open("POST", "http://localhost:8080/login", true);  
        xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");  
        xhr.onreadystatechange = () => {  
            if (xhr.readyState === 4) {  
                const response = JSON.parse(xhr.responseText); // 解析 JSON 回應  
                if (xhr.status === 200) {  
                    if (response.code === 0) { // 檢查返回的 code   
                        sys.localStorage.setItem('username', response.data.Username);  
                        sys.localStorage.setItem('money', response.data.Money);  

                       
                        console.log("Login success: " + response.message); 
                        console.log("Login data: " + response.data.Money); 
                        director.loadScene("A"); // 切換到場景 A  

                    } else {  
                        console.log("Login failed: " + response.message);  
                        this.usernameInput.string = ""; // 清空帳號  
                        this.passwordInput.string = ""; // 清空密碼  
                    }  
                } else {  
                    console.log("Login request failed: " + response.message);  
                    this.usernameInput.string = ""; // 清空帳號  
                    this.passwordInput.string = ""; // 清空密碼  
                }  
            }  
        };  
        xhr.send(JSON.stringify(loginData));  
    }  

    
}