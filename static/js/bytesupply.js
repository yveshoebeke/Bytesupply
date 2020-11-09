$(document).ready(function() {
    setInterval(function() {
        var time = new Date();
        $("#doy").text(time.toDateString());
        $("#tod").text(time.toLocaleTimeString());
    }, 1000);

    $(".lang-es").hide();
    
    $(".logo").mouseenter(function(){
        $("#logoExplain").text($(this).attr("alt") + " -> " + $(this).attr("title"));
    }).mouseout(function() {
        $("#logoExplain").text("");
    });

    $("#contactForm").submit(function() {
        //event.preventDefault();

        var name = $("#contactName").val();
        var company = $("#contactCompany").val();
        var phone = $("#contactPhone").val();
        var email = $("#contactEmail").val();
        var URL = $("#contactURL").val();
        var message = $("#contactMessage").val();

        grecaptcha.ready(function() {
            grecaptcha.execute('6LePhd8ZAAAAAATEfk8aVTl4VJXj8R8OKgStdAzI', {action: 'submit'}).then(function(token) {
                console.log("token: " + token); // Add your logic to submit to your backend server here.
                $("contactForm").prepend('<input type="hidden" name="g-recaptcha-response" value="' + token + '">');

                $.post("/contactus",{
                                    validEntry: "false",
                                    contactName: name,
                                    contactCompanay: company,
                                    contactPhone: phone,
                                    contactEmail: email,
                                    contactURL: URL,
                                    contactMessage: message, 
                                    token: token}, 
                                    function(result) {
                                        console.log(result);
                                        /*
                                        if(result.success) {
                                                alert('Thanks for posting comment.')
                                        } else {
                                                alert('You are spammer ! Get the @$%K out.')
                                        }
                                        */
                                    });
            });
        });    
    });

});

function validateMessageEntry(id) {
    Id = "#" + id;
    errId = Id + "Help";
    validEntry = false;
    $(errId).removeClass("d-block").addClass("d-none");

    if (($(Id).hasClass("mandatory") && $(Id).val().length == 0) || ($(Id).val().length > 256)) {
        validEntry = false;
    } else {
        if (id == "contactMessage") {
            type = "text";
        } else {
            type = $(Id).attr("type");
        }

        switch(type) {
            case 'tel':
                if ($(Id).val().length == 0){
                    break;
                } else if ($(Id).val().length > 13){
                    validEntry = false;
                    break;
                } else {
                    validEntry = (/^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$/.test($(Id).val()));
                    break;
                }
            case 'email':
                validEntry = (/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+$/.test($(Id).val()));
                break;
            case 'url':
                validEntry = (/^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$/gm.test($(Id).val()));
                break;
            case 'text':
            default:
                validEntry = (/^[a-zA-Z0-9 \-\,\.&]+$/.test($(Id).val()));
                break;
        }
 
        if (!validEntry) {
            $("#validEntry").val("false");
            $(errId).removeClass("d-none").addClass("d-block");
            $(Id).focus();
        } else {
            $("#validEntry").val("true");
        }
    }
}

function setLang(langId) {
    if (langId == "lang-es") {
        oldLangId = "lang-en"
    } else {
        oldLangId = "lang-es"
    }

    $("#"+langId).attr("src", "/static/img/" + oldLangId + ".png");
    $("#"+langId).attr("id", oldLangId);

    $("."+oldLangId).hide(function() {$("."+langId).show()});
}

/*
*/
