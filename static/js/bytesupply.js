$(document).ready(function() {
    if(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4))){
        document.querySelector("body").classList.remove("desktopBody")
        document.querySelector("body").classList.add("mobileBody")
    }

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

    // Set session space
    // name -> username
    window.sessionStorage;




    // Remove after revamping --> see register/login logic
    $("#contactForm").submit(function() {
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
    // ^^^ Remove after revamping --> see register/login logic


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

let Register = class {
    constructor(userName, realName, title, loginTime, lastLogin) {
        this.userName = userName
        this.realName = realName
        this.title = title
        this.loginTime = loginTime
        this.lastLogin = lastLogin
    }

    // Register user in sessionStorage
    setSessionLoginData(){
        sessionStorage.setItem("userName", this.userName)
        sessionStorage.setItem("realName", this.realName)
        sessionStorage.setItem("title", this.title)
        sessionStorage.setItem("loginTime", this.loginTime)
        sessionStorage.setItem("lastLogin", this.lastLogin)
    }
}

// Get register data element
function getLoginElement(key) {
    return sessionStorage.getItem(key);
}

// Remove registration
function unsetLoginData(){
    sessionStorage.clear();
}
/*
*/
