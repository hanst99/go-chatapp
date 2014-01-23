//reload chat every 5 seconds, and when a message was posted by ourselves
$(function() {
    function reloadChat() {
        $('#chatbox').load('/ #chatbox>li',function() {})
    }
    window.setInterval(reloadChat,5000)
    $("#post_button").click(function() {
        $.ajax({
            type: 'POST',
            url: '/post_message',
            data: $("#message").val(),
            success: function(_data) {
                reloadChat()
            }
        })
    })
})

