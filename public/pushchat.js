$(function() {
    var chatboxIsScrolledDown = function() {
        var chatbox = $('#chatbox')
        return chatbox.prop('scrollHeight') - chatbox.scrollTop() == chatbox.outerHeight()
    }
    var chatBoxScrollDown = function () {
        var chatbox = $('#chatbox')
        chatbox.scrollTop(chatbox.prop('scrollHeight'))
    }
    var loc = window.location
    if(!("WebSocket" in window)) {
        $('<li><p class="from">System</p>' +
            '<p class="message">You need a recent browser with support for WebSockets for this page to work!</p>'+
            '</li>').appendTo('#chatbox')
    } else {
        var url = 'ws://'+loc.hostname+':'+loc.port+'/chat'
        var ws = new WebSocket(url)
        ws.onmessage = function(smsg) {
           var scrolledDown = chatboxIsScrolledDown()
           var msg = JSON.parse(smsg.data) 
           $('<li><p class="from">['+msg.From+']</p><p class="message">'+msg.Content+'</li>').appendTo('#chatbox')
           if(scrolledDown) chatBoxScrollDown()
        }
        $('#message_form').submit(function() {
            var msg = JSON.stringify({From: $('#name').val(), Content: $('#message').val()})
            chatBoxScrollDown()
            ws.send(msg)
            return false
        })
    }
})
