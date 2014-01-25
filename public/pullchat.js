//reload chat every 5 seconds, and when a message was posted by ourselves
$(function() {
    var reloadChat = function(alwaysScroll) {
        var chatbox = $('#chatbox')
        //scroll down after loading if we're told to do so or if the chatbox was completely scrolled down already
        var isScrolledDown = chatbox.prop('scrollHeight') - chatbox.scrollTop() == chatbox.outerHeight()
        var scrollAfterDone = alwaysScroll || isScrolledDown
        chatbox.load('/ #chatbox>li',function() {
            if(scrollAfterDone) {
                chatbox.scrollTop(chatbox.prop('scrollHeight'))
            }
        })
        
    }
    window.setInterval(function(){reloadChat(false)},5000)
    $("#message_form").submit(function() {
        $.ajax({
            type: 'POST',
            url: '/post_message',
            data: $("#message").val(),
            success: function(_data) {
                reloadChat(true)
            }
        })
        //don't reload page - causes
        //warnings in recent chrome versions.
        //This appears to be caused by jQuery internals
        //and isn't likely to cause any actual problems now or in the future
        return false
    })
})

