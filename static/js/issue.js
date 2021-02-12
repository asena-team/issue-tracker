const formSubmitEvent = (e) => {
  const isNameBlank = document.getElementById("name").value === "";
  const isEmailBlank = document.getElementById("email").value === "";
  const isTitleBlank = document.getElementById("title").value === "";
  const isDescBlank = document.getElementById("desc").value === "";

  const inputInfo = [
    {
      name: "name",
      blank: isNameBlank
    }, {
      name: "mail",
      blank: isEmailBlank,
    }, {
      name: "title",
      blank: isTitleBlank
    }, {
      name: "desc",
      blank: isDescBlank
    },
  ];

  for (let i = 0; i < inputInfo.length; i++) {
    if (inputInfo[i].blank === true) {
      const info = inputInfo[i].name;
      infoMessages(info)
    }
  }

  var typeDefaultOption = document.getElementById("type").options[0].value;
  var priotiryDefaultOption = document.getElementById("priority").options[0]
    .value;

  $("#name").val("");
  $("#email").val("");
  $("#title").val("");
  $("#type").val(typeDefaultOption);
  $("#priority").val(priotiryDefaultOption);
  $("#desc").val("");

  infoMessages("success")

  e.preventDefault();
}

const infoMessages = (inputName) => {
  if (inputName === "success") {
    $("#success-modal").show();
    $(".name-info-message").hide()
    $(".mail-info-message").hide()
    $(".title-info-message").hide()
    $(".desc-info-message").hide()
    setTimeout(() => {
      $("#check").attr("checked", true);
    }, 1500);

    setTimeout(() => {
      $("#check").attr("checked", false);
      $("#success-modal").hide();
    }, 2200);
  } else {
    $(`.${inputName}-info-message`).show()
  }
};

const form = document.getElementById("issue-form");
form.addEventListener("submit", formSubmitEvent);
