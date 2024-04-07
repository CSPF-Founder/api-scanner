export function ready(fn) {
  if (document.readyState !== 'loading') {
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

var alert_close_button = document.getElementsByClassName('alert-close-button');

for (let i = 0; i < close.length; i++) {
  alert_close_button[i].onclick = function () {
    var div = this.parentElement;
    div.style.opacity = '0';
    setTimeout(function () {
      div.style.display = 'none';
    }, 30);
  };
}

const appMsgBox = document.querySelector('#app-msg-box');

let appMsgBoxTitle = null;
let appMsgBoxClose = null;
let appMsgBoxHeader = null;
let appMsgBoxBody = null;

if (appMsgBox) {
  appMsgBoxTitle = appMsgBox.querySelector('.modal-title');
  appMsgBoxBody = appMsgBox.querySelector('.modal-body');
  appMsgBoxHeader = appMsgBox.querySelector('.modal-header');
  appMsgBoxClose = appMsgBox.querySelector('.btn-close');
}

document.addEventListener('DOMContentLoaded', function () {
  if (appMsgBoxClose) {
    appMsgBoxClose.addEventListener('click', function () {
      setTimeout(function () {

        appMsgBoxTitle.textContent = '';
        appMsgBoxBody.textContent = '';
        
        appMsgBox.style.display = 'none';
      }, 300);
    });
  }

  // When the user clicks anywhere outside of the modal, close it
  window.onclick = function (event) {
    if (event.target === appMsgBox) {
      setTimeout(function () {

        appMsgBoxTitle.textContent = '';
        appMsgBoxBody.textContent = '';
        
        appMsgBox.style.display = 'none';
      }, 300);
    }
  };
});


export function showModalMessage(message, title) {
  setTimeout(function () {
      appMsgBoxTitle.textContent = title;
      appMsgBoxBody.textContent = message;
      appMsgBox.style.display = 'block';
  }, 600); // 600 milliseconds wait before displaying the box
}

export function showSuccess(message, title) {
  hideLoadingBox();

  if (!title) {
    title = 'Success';
  }

  appMsgBoxHeader.classList.remove('bg-primary', 'bg-warning');
  appMsgBoxHeader.classList.add('bg-success');

  showModalMessage(message, title);
}

export function showInfo(message, title) {
  if (!title) {
    title = 'Info';
  }

  appMsgBoxHeader.classList.remove('bg-warning', 'bg-success');
  appMsgBoxHeader.classList.add('bg-primary');

  showModalMessage(message, title);
}

export function showError(message, title) {
  if (!title) {
    title = 'Error';
  }

  if (message instanceof Array) {
    message = message.map(msg => '\u2022 ' + msg).join('\n');
  }

  appMsgBoxHeader.classList.remove('bg-primary', 'bg-success');
  appMsgBoxHeader.classList.add('bg-warning');

  showModalMessage(message, title);
}

export function smallMessageBox(msg) {
  let messageBox = document.getElementById('message_box');

  if (!messageBox) {
    // Constructing the HTML string
    let msgHTML = `<div class="info" id="message_box"><table class="mx-auto"><tr><td id="message">${msg}</td>`;
    msgHTML += `<td></td></tr></table></div>`;

    // Appending the new element to the '#content' element
    const content = document.getElementById('content');
    if (content) {
      content.insertAdjacentHTML('beforeend', msgHTML);
    }
  } else {
    // Update the message if the element already exists
    const message = document.getElementById('message');
    if (message) {
      message.innerHTML = msg;
    }
  }

  // Displaying the message box
  if (messageBox) {
    messageBox.style.display = 'block';
  }
}

export function hideSmallMessageBox() {
  setTimeout(function () {
    const messageBox = document.getElementById('message_box');
    if (messageBox) {
      messageBox.style.display = 'none';
    }
  }, 100);
}

export function loadingBox() {
  smallMessageBox(
    '<b>Sending Request .. <i class="fa fa-spinner fa-pulse fa-2x fa-fw"></i></b>',
    -1
  );
}

export function hideLoadingBox() {
  hideSmallMessageBox();
}

export function createHiddenInputElement(name, value) {
  var hiddenField = document.createElement('input');
  hiddenField.setAttribute('type', 'hidden');
  hiddenField.setAttribute('name', name);
  hiddenField.setAttribute('value', value);
  return hiddenField;
}

export function virtualFormSubmit(path, params, method) {
  //Reference: http://ctrlq.org/code/19233-submit-forms-with-javascript
  method = method || 'post';
  var form = document.createElement('form');
  form.setAttribute('method', method);
  form.setAttribute('action', path);
  for (var name in params) {
    if (params.hasOwnProperty(name)) {
      var value = params[name];
      if (value && value instanceof Array) {
        for (var i = 0; i < value.length; i++) {
          hiddenField = createHiddenInputElement(name, value[i]);
          form.appendChild(hiddenField);
        }
      } else {
        hiddenField = createHiddenInputElement(name, value);
        form.appendChild(hiddenField);
      }
    }
  }
  document.body.appendChild(form);
  form.submit();
}

export function resetInputForm(form_id) {
  if (form_id !== undefined) {
    // Query the form by ID and reset it if it exists
    const form = document.querySelector(form_id);
    if (form) {
      form.reset();
    }
  } else {
    // Query all forms with the class '.input_form' and reset each
    document.querySelectorAll('.input_form').forEach(form => {
      form.reset();
    });
  }
}

export function redirect(url, msg) {
  showError(msg);
  setTimeout(function () {
    window.location.href = url;
  }, 5000);
}

export function redirectToLogin(url) {
  if (
    window.confirm(
      'You have no longer logged in, you will be redirected to Login page..'
    )
  ) {
    window.location.href = url;
  }
}

ready(function() {

  resetInputForm();

  // Handling the close button on alert boxes
  const alertBoxCloseButtons = document.getElementsByClassName('alert-box-close');
  for (let i = 0; i < alertBoxCloseButtons.length; i++) {
    alertBoxCloseButtons[i].addEventListener('click', function () {
      const div = this.parentElement;
      div.style.opacity = '0';
      setTimeout(function () {
        div.style.display = 'none';
      }, 600);
    });
  }

  // Handling the 'select all' functionality in tables
  const selectAllRowsCheckboxes = document.querySelectorAll('.select-all-rows');
  selectAllRowsCheckboxes.forEach(function (checkbox) {
    checkbox.addEventListener('click', function (e) {
      const currentTable = this.closest('table');
      const checkboxes = currentTable.querySelectorAll("input[name='selected_rows[]']");
      checkboxes.forEach(function (cb) {
        cb.checked = checkbox.checked;
      });
    });
  });
});

document.addEventListener('DOMContentLoaded', function () {
  // Append CSRF token as hidden input to all forms
  document.querySelectorAll('form').forEach(function (form) {
    let csrfInput = document.createElement('input');
    csrfInput.setAttribute('type', 'hidden');
    csrfInput.setAttribute('name', 'csrf_token');
    csrfInput.setAttribute('value', CSRF_TOKEN);
    form.appendChild(csrfInput);
  });
});

  // Function to include CSRF token in Fetch requests
export function requestWithCSRFToken(url, options = {}) {
  // Set default method to GET
  options.method = options.method || 'GET';

  // Set headers if not provided
  if (!options.headers) {
    options.headers = {};
  }

  // Add CSRF token to headers for all requests
  options.headers['X-CSRF-Token'] = CSRF_TOKEN;
  options.headers['X-Requested-With'] = 'XMLHttpRequest';

  // Add CSRF token to POST request body
  if (options.method.toUpperCase() === 'POST') {
    let data = options.body || new FormData();
    if (typeof data === 'string' && !data.includes(CSRF_NAME)) {
      data += `&${CSRF_NAME}=${CSRF_TOKEN}`;
    } else if (data instanceof FormData) {
      if (!data.has(CSRF_NAME)) {
        data.append(CSRF_NAME, CSRF_TOKEN);
      }
    }
    options.body = data;
  }

  // Perform the fetch request
  return fetch(url, options);
}

// Use requestWithCSRFToken instead of fetch for all requests
// Example:
// requestWithCSRFToken('/some-url', { method: 'POST', body: 'data' });


ready(function () {
  // is it required?
    let url = window.location.href;

    // Query all anchor tags inside '#sidebar'
    let sidebarLinks = document.querySelectorAll('#sidebar a');

    sidebarLinks.forEach(link => {
        if (link.href === url) {
            // Add 'active' class to the parent of the link
            link.parentElement.classList.add('active');
        }
    });
});
