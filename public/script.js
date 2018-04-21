//// websocket start
var player_name;
var chess_map;
var selected_chess;
var chess_term;
var user_id;
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
};
ws.onopen = function() {
  $("#connect_status").html("已连接");
};
ws.onclose = function (evt) {
  $("#connect_status").html("连接已关闭");
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
        str = "[" + x + "," + y + "]"
        if(map[str]==""){ // 无棋子坐标
          $("#chess_map").append("<div class='block empty' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
        }else{
          if(map[str][0]=="R"){
            $("#chess_map").append("<div class='block r-chess' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
          }else{
            $("#chess_map").append("<div class='block b-chess' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
          }
          
        }
      }
      $("#chess_map").append("<br>");
    }
  }else{
    for(var y = 9;y >=0; y--){
      for(var x = 0;x <=8; x++){
        str = "[" + x + "," + y + "]"
        if(map[str]==""){ // 无棋子坐标
          $("#chess_map").append("<div class='block empty' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
        }else{
          if(map[str][0]=="R"){
            $("#chess_map").append("<div class='block r-chess' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
          }else{
            $("#chess_map").append("<div class='block b-chess' onclick='chess_click(this,\""+map[str]+"\",\""+str+"\")'>"+map[str]+"<br>"+str+"</div>");
          }
          
        }
      }
      $("#chess_map").append("<br>");
    }
  }
  
};
/////

var showMsgToPanel = function(){
  if(chess_term == 1){
    $("#next_term").html("蓝棋");
  }else{
    $("#next_term").html("红棋");
  }

  if("R"==current_color){
    $("#current_color").html("红棋");
  }else if ("B"==current_color){
    $("#current_color").html("蓝棋");
  }
  
};

var move = function(user_id,player_name,left,top,callback){
  $("#"+player_name).animate({left: left+"px",top: top+"px"},"fast");
  if(callback!=null){
    callback(player_name,origin,target);
  }
};


function send_msg(user_id,player_name,operate,origin,target){
  ws.send(JSON.stringify({user_id,player_name,operate,origin,target}));
}
//// block controller end

var set_player_name_default = function(){
  player_name = "cheng";
  console.log("init default method")
  if(player_name){
      $("#hi_title").html(player_name);
      // send_msg(user_id,player_name,"INIT_GAME","","");
      send_msg(user_id,player_name,"RELOAD_GAME","","");
  }
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
    $("#selected_chess").html(chess+point);
    selected_chess = point;
    $(self).css("background-color","red");
  }
};

var clear = function (){ //清除焦点
  $("#selected_chess").html("");
  selected_chess = undefined;
  $(".block").css("background-color","rgba(255, 255, 0, 0.166)");
}

$().ready(function(){
  user_id = getCookie('chinese_chess_user_id');
  setTimeout("set_player_name_default()",500);
  $("#restart").click(function(){
    send_msg(user_id,player_name,"RELOAD_GAME","","");
  });
  $("#clear_btn").click(clear);
});

function getCookie(name){
  var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");
  if(arr=document.cookie.match(reg))
    return unescape(arr[2]);
  else
    return null;
}
// save player name end