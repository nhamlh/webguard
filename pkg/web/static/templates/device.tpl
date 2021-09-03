<html>
  {{ template "header" }}
  <body>
    <div class="container">
      <br>
      {{ template "display_errors" . }}
      <br>

      <div class="device-form" align="center">
        <form action="/new_device" method="POST">

          <div class="field">
            <label class="label">Name</label>
            <div class="control">
              <input class="input" type="text" name="name" placeholder="My Device">
            </div>
          </div>

          <div class="field is-grouped">
            <div class="control">
              <button class="button is-primary">Submit</button>
            </div>
            <div class="control">
              <button class="button is-link is-light">Cancel</button>
            </div>
          </div>
        </form>
      </div>

    </div>
  </body>
</html>
