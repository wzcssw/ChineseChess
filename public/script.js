//// websocket start
var player_name;
var chess_map;
var selected_chess;
var chess_term;
var user_id = null;
var current_color;
var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
  uri = 'wss:';
}
uri += '//' + loc.host;
uri += loc.pathname + 'ws';

ws = new WebSocket(uri);
if(window.WebSocket != undefined) {
  $("#connect_status").html("浏览器不支持websockets");
  $("#connect_status").css("color","red");
};
ws.onopen = function() {
  $("#connect_status").html("已连接");
  $("#connect_status").css("color","green");
};
ws.onclose = function (evt) {
  $("#connect_status").html("连接已关闭");
  $("#connect_status").css("color","red");
}; 
ws.onmessage = function(evt) {
  if(evt){
    var msg = evt.data;
    console.log("CC",msg)
    obj = JSON.parse(msg);
    if(obj){
      clear();
      msg_process(obj);
    }else{
      console.log('数据错误',evt)
      alert('数据错误');
    }
  }else{
    alert('服务器错误');
  }
};

this.send = function (message, callback) {  
  this.waitForConnection(function () {  
      ws.send(message);  
      if (typeof callback !== 'undefined') {  
        callback();  
      }  
  }, 1000);  
};

this.waitForConnection = function (callback, interval) {  
  if (ws.readyState === 1) {  
      callback();  
  } else {  
      var that = this;  
      // optional: implement backoff for interval here  
      setTimeout(function () {  
          that.waitForConnection(callback, interval);  
      }, interval);  
  }  
};  
//// websocket end

/////
var msg_process = function(chess){
  var map = chess.map;
  chess_term = chess.term;
  if (chess.r_user_id == user_id){ // 选择当前阵营
    current_color = "R"
  }else{
    current_color = "B"
  }
  showMsgToPanel();
  render_chess(map);
};

var render_chess = function(map){
  $("#chess_map").html("");
  if(current_color=="B"){
    for(var y = 0;y <=9; y++){
      for(var x = 8;x >=0; x--){
        appent_chess_dom(map,x,y);
      }
      $("#chess_map").append("<br>");
    }
  }else{
    for(var y = 9;y >=0; y--){
      for(var x = 0;x <=8; x++){
        appent_chess_dom(map,x,y);
      }
      $("#chess_map").append("<br>");
    }
  }
  
};

var appent_chess_dom = function(map,x,y){
  str = "[" + x + "," + y + "]"
  cls = "empty";
  if (map[str][0]=="R"){
    cls = "r-chess"
  }else if(map[str][0]=="B"){
    cls = "b-chess"
  }
  $("#chess_map").append("<div class=\"block "+cls+" \" onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
}
/////

var showMsgToPanel = function(){
  if(chess_term == 1){
    $("#next_term").html("蓝棋");
    $("#next_term").css("color","blue");
  }else{
    $("#next_term").html("红棋");
    $("#next_term").css("color","red");
  }

  if("R"==current_color){
    $("#current_color").html("红棋");
    $("#current_color").css("color","red");
  }else if ("B"==current_color){
    $("#current_color").html("蓝棋");
    $("#current_color").css("color","blue");
  }
  
};

var move = function(user_id,player_name,left,top,callback){
  $("#"+player_name).animate({left: left+"px",top: top+"px"},"fast");
  if(callback!=null){
    callback(player_name,origin,target);
  }
};


function send_msg(user_id,player_name,operate,origin,target){
  this.send(JSON.stringify({user_id,player_name,operate,origin,target}));
}
//// block controller end

var connect_to_server = function(){
  send_msg(user_id,player_name,"RELOAD_GAME","","");
};


// 棋子按下事件
var chess_click = function(self,chess,point){
  // 判断是否当前用户是否可以出棋
  var cc = (current_color=="R" ? 0:1);
  if(cc!=chess_term){
    return 
  }

  if(selected_chess){ // 如果之前已选中棋子
    send_msg(user_id,player_name,"MOVE",selected_chess,point);
  }else{// 如果未选中棋子
    if (chess[0]!=current_color){ //判断是否是自己的棋子
      return
    }
    if (chess==""){
      return 
    }
    selected_chess = point;
    $(self).css("background-color","red");
  }
};

var clear = function (){ //清除焦点
  $("#selected_chess").html("");
  selected_chess = undefined;
  $(".block").css("background-color","rgba(255, 255, 0, 0.166)");
}

var save_username = function(){
  var username = $("#username").val();
  setCookie("chinese_chess_username",username);
}

$().ready(function(){
  user_id = getCookie('chinese_chess_user_id');
  player_name = getCookie('chinese_chess_username');
  $("#hi_title").html(player_name);
  setTimeout("connect_to_server()",500);
});

function getCookie(name){
  var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");
  if(arr=document.cookie.match(reg))
    return unescape(arr[2]);
  else
    return null;
}

function setCookie(c_name,value,expiredays){
  var exdate=new Date()
  exdate.setDate(exdate.getDate()+expiredays)
  document.cookie=c_name+ "=" +escape(value)+
  ((expiredays==null) ? "" : ";expires="+exdate.toGMTString())
}