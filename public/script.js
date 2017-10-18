//// websocket start
var player_name;
var chess_map;
var selected_chess;
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
    console.log(msg);
    obj = JSON.parse(msg);
    if(obj){
      clear();
      console.log(msg);
      msg_process(obj);
    }else{
      alert('数据错误');
    }
  }else{
    alert('服务器错误');
  }
};
//// websocket end

/////
var msg_process = function(msg){
  $("#chess_map").html("");
  for(var y = 9;y >=0; y--){
    for(var x = 0;x <=8; x++){
      str = "[" + x + "," + y + "]";
      if(msg[str]==""){
        msg[str]="x";
      }
      // $("#chess_map").append(msg[str]);
      $("#chess_map").append("<div class='block' onclick='chess_click(this,\""+msg[str]+"\",\""+str+"\")'>"+msg[str]+"<br>"+str+"</div>");
    }
    $("#chess_map").append("<br>");
  }
};
/////

// //// block controller start
// document.onkeydown = function(event){
//     var e = event || window.event || arguments.callee.caller.arguments[0];
//     var offset_left = $("#"+player_name).position().left;
//     var offset_top = $("#"+player_name).position().top;
//     if(e && e.keyCode==37){ // Left
//       move(player_name,offset_left-20,offset_top,send_msg);
//     };
//     if(e && e.keyCode==38){ // 按 Up 
//       move(player_name,offset_left,offset_top-20,send_msg);
//     };
//     if(e && e.keyCode==39){ // 按 Right 
//       move(player_name,offset_left+20,offset_top,send_msg);
//     };
//     if(e && e.keyCode==40){ // Down
//       move(player_name,offset_left,offset_top+20,send_msg);
//     };
// };

var move = function(player_name,left,top,callback){
  $("#"+player_name).animate({left: left+"px",top: top+"px"},"fast");
  if(callback!=null){
    callback(player_name,origin,target);
  }
};

function send_msg(player_name,operate,origin,target){
  ws.send(JSON.stringify({player_name,operate,origin,target}));
}
//// block controller end

// save player name start
var set_player_name = function(){
  if(!player_name){
      var input = $("#player_name_input").val();
      if(input == ""){
        alert("alert");
      }else{
        player_name = input;
        $("#hi_title").html(player_name);
        add_block(player_name,0,0);
        // ws.send(JSON.stringify({player_name,left,top}));
        send_msg(player_name,0,0);
      }
  }
};

var set_player_name_default = function(){
  player_name = "cheng";
  console.log("init default method")
  if(player_name){
      $("#hi_title").html(player_name);
      add_block(player_name,0,0);
      // ws.send(JSON.stringify({player_name,left,top}));
      send_msg(player_name,"INIT_GAME","","");
  }
};

var add_block = function(name,top,left){
  // $("body").append("<div style='left: "+left+"px;top: "+top+"px' class='block "+name+"' id='"+name+"'>"+name+"</div>");
};

// 棋子按下事件
var chess_click = function(self,chess,point){
  if(selected_chess){ // 如果之前已选中棋子
    send_msg(player_name,"MOVE",selected_chess,point);
  }else{// 如果未选中棋子
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
  setTimeout("set_player_name_default()",500);
  $("#ok_btn").click(function(){
    set_player_name();
  });
  $("#restart").click(function(){
    send_msg(player_name,"INIT_GAME","","");
  });
  $("#clear_btn").click(clear);
});
// save player name end