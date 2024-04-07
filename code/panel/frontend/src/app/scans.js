import {
  redirectToLogin,
  showError,
  showSuccess,
  resetInputForm,
  loadingBox,
  hideLoadingBox,
  requestWithCSRFToken,
  ready,
} from './main.js';

import 'datatables.net';
import 'datatables.net-bs4';

// Only Jquery Dependency
$(document).ready(function () {
  $('#scan-list').DataTable({
    searching: false,
  });

});

ready(function () {
  const addNewHttpHeaders =   document.querySelector('.add_new_http_headers');
  const httpHeadersDiv =  document.getElementById("http-headers-div");
  const httpHeaders =   document.querySelector('.http_headers');
  let x = 1;

  if (!addNewHttpHeaders || !httpHeaders || !httpHeadersDiv) {
    return;
  }

  addNewHttpHeaders.addEventListener('click', function () {
    x++;
    httpHeaders.insertAdjacentHTML('beforeend', '<div class="pt-4 http-headers-div"> \
                                <div class="row"> \
                                    <div class="col-lg-2 col-sm-0"></div> \
                                    <div class="col-lg-2 col-sm-12"> \
                                        <label>Header Name</label> \
                                    </div> \
                                    <div class="col-lg-6 colo-sm-12"> \
                                        <select class="form-control header_name" name="http_headers[][header_name]"> \
                                            <option value="authorization">Authorization</option> \
                                            <option value="x_api_key">x-api-key</option> \
                                            <option value="custom">Custom</option> \
                                        </select> \
                                    </div> \
                                    <a class="circle-icon bg-red tooltip-bs-style1 remove_button" href="javascript:void(0);" style="margin-top:6px; width:21px; height:21px;"> \
                                        <i class="fa fa-minus text-white" style="margin-left:-6px;"> \
                                            <span class="tooltip-bs-style1-text"> Remove </span> \
                                        </i> \
                                    </a>  \
                                </div> \
                                <div class="row div-custom-header-name" style="display:none;"> \
                                    <div class="col-lg-2 col-sm-0"></div> \
                                    <div class="col-lg-2 col-sm-12"></div> \
                                    <div class="col-lg-6 col-sm-12"> \
                                        <input type="text" class="form-control sharpedge" name="http_headers[][custom_header_name]" placeholder="Custom Header Name"> \
                                    </div> \
                                </div> \
                                <div class="row"> \
                                    <div class="col-lg-2 col-sm-0"> </div> \
                                    <div class="col-lg-2 col-sm-12"> \
                                        <label>Header Value</label> \
                                    </div> \
                                    <div class="col-lg-6 col-sm-8"> \
                                        <input type="text" class="form-control sharpedge" name="http_headers[][header_value]"> \
                                    </div> \
                                </div> \
                            </div>');
    
  });

  httpHeadersDiv.addEventListener('click', function (e) {

    let targetElement = e.target;

    // Traverse up the DOM tree if the clicked element is not the button itself
    while (targetElement != null && !targetElement.classList.contains('remove_button')) {
        targetElement = targetElement.parentElement;
    }

    // If no element with the 'remove_button' class was found in the ancestry, exit the function
    if (targetElement == null) return;

    e.preventDefault();

    var closestDiv = e.target.closest('div.http-headers-div');
    if (closestDiv) {
        closestDiv.remove();
    }

  });

  httpHeadersDiv.addEventListener('change', function (e) {
    let targetElement = e.target;

    while (targetElement != null && !targetElement.classList.contains('header_name')) {
        targetElement = targetElement.parentElement;
    }
  
    if (targetElement == null) return;

    let parentDiv = e.target.closest('.http-headers-div');

    let header_name = e.target.value;

    let customHeaderNameDiv = parentDiv.querySelector('.div-custom-header-name');

     if (header_name === 'custom') {
        customHeaderNameDiv.style.display = '';
    } else {
        customHeaderNameDiv.style.display = 'none';
    }

  });

});


// Delete scan
ready(function () {
  const scanList = document.getElementById('scan-list');

  if (!scanList) {
    return;
  }

  scanList.addEventListener('click', function (e) {
    let targetElement = e.target;

    // Traverse up the DOM tree if the clicked element is not the button itself
    while (targetElement != null && !targetElement.classList.contains('delete-scan')) {
        targetElement = targetElement.parentElement;
    }

    // If no element with the 'delete-scan' class was found in the ancestry, exit the function
    if (targetElement == null) return;

    e.preventDefault();
    if (!confirm('Are you sure want to delete the scan ?')) {
      return;
    }

    let row = e.target.closest('tr');
    let scan_id = row.getAttribute('data-id'); // Adjust if data-id is stored differently
    let url = '/scans/' + scan_id;

    loadingBox(); 
    requestWithCSRFToken(url, {
        method: 'DELETE',
      })
      .then(response => response.json().then(data => ({ ok: response.ok, data })))
      .then(({ ok, data }) => {
        hideLoadingBox();

        if (!ok) {
          throw new Error(data.error || 'Error occurred');
        }

        if (data.redirect) {
          redirectToLogin(data.redirect);
        } else if (data.success) {
          showSuccess(data.success);
          row.remove();
        }
      })
      .catch(error => {
        hideLoadingBox();
        showError(error.message);
      });

  });

});


// add scan
ready(function () {
  const addForm = document.getElementById('add-scan-form');
  if (!addForm) {
    return;
  }

  const addButton = document.getElementById('add-scan-btn');

  // Add Scan Form Submit
  addForm.addEventListener('submit', function (event) {
    addButton.disabled = true;

    event.preventDefault();
    loadingBox();

    const formData = new FormData(addForm);

    requestWithCSRFToken('/scans', {
      method: 'POST',
      body: formData,
    })
    .then(response => response.json().then(data => ({ ok: response.ok, data })))
    .then(({ ok, data }) => {
      hideLoadingBox();

      if (!ok) {
        throw new Error(data.error || 'Error occurred');
      }

      if (data.success) {
        resetInputForm('#add-scan-form');
        showSuccess(data.success);
      } else if (data.redirect) {
        redirectToLogin(data.redirect);
      }

      addButton.disabled = false;
    })
    .catch(error => {
      hideLoadingBox();
      showError(error.message);
      addButton.disabled = false;
    });

  });

});


