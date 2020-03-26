// // (function($) {
// //     $.fn.postRegister = function() {
// //         alert('hello world');
// //         return this;
// //     };
// // })(jQuery);

// function postRegister() {
//     console.log("postRegister func is triggered")
//     return p1 * p2; // The function returns the product of p1 and p2
// }

$(document).ready(function() {
    $("form_submit").click(function(e) {
        e.preventDefault();
        alert("button");
        console.log($(this).val())

    });

    $("#form_submit").click(function(e) {
        e.preventDefault();
        console.log($(this).val())
        var name, email, password, passwordConfirmation, phoneNo
        $('form input, form select').each(
            function(index) {
                var input = $(this);
                // console.log('Name: ' + input.attr('name') + ' Value: ' + input.val());
                if (input.attr('name') == 'fullName') {
                    name = input.val()
                } else if (input.attr('name') == "email") {
                    email = input.val()
                } else if (input.attr('name') == "password") {
                    password = input.val()
                } else if (input.attr('name') == "passwordConfirmation") {
                    passwordConfirmation = input.val()
                } else if (input.attr('name') == "phoneNo") {
                    phoneNo = input.val()
                }
            }
        );

        $.ajax({
            type: "POST",
            url: "/api/register",
            dataType: 'json',
            contentType: 'application/json',
            data: JSON.stringify({
                name: name,
                email: email,
                phone_no: "8090898778",
                password: password,

            }),
            success: function(result) {
                alert('ok');
            },
            error: function(result) {
                alert('error');
            }
        });



    });
});