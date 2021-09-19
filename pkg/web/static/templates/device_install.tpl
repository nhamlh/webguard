<html>
  {{ template "header" . }}
  <body>
    <div class="container is-fluid">

      {{ template "display_errors" . }}

    <section id="install" class="section">
      <div class="is-divider" data-content="Step one"></div>
      <div id="step-one" class="columns is-centered">
        <div class="column has-text-centered">
          <div class="content">
            <h3>Download and install VPN client for your device</h3>
                <a class="button is-link is-outlined" href="https://itunes.apple.com/us/app/wireguard/id1451685025?ls=1&mt=12">
                  <span class="icon">
                    <i class="fab fa-apple"></i>
                  </span>
                  <span>MacOS</span>
                </a>

                <a class="button is-link is-outlined" href="https://download.wireguard.com/windows-client/wireguard-installer.exe">
                  <span class="icon">
                    <i class="fab fa-windows"></i>
                  </span>
                  <span>Windows</span>
                </a>

                <a class="button is-link is-outlined" href="https://www.wireguard.com/install/#ubuntu-module-tools">
                  <span class="icon">
                    <i class="fab fa-linux"></i>
                  </span>
                  <span>Linux</span>
                </a>

                <a class="button is-link is-outlined" href="https://play.google.com/store/apps/details?id=com.wireguard.android">
                  <span class="icon">
                    <i class="fab fa-android"></i>
                  </span>
                  <span>Android</span>
                </a>

                <a class="button is-link is-outlined" href="https://itunes.apple.com/us/app/wireguard/id1441195209?ls=1&mt=8">
                  <span class="icon">
                    <i class="fab fa-app-store-ios"></i>
                  </span>
                  <span>iOS/ipadOS</span>
                </a>

          </div>
        </div>
      </div>
      <div class="is-divider" data-content="Step two"></div>
      <div id="step-two" class="columns is-centered">
        <div class="column has-text-centered">
          <div class="content">
            <h3>For PC (Mac/Linux/Windows)</h3>
            <p><span>Click to </span>
              <span>
              <a class="button is-primary is-small" href="{{ $.download_url }}" download>
                <span class="icon is-small">
                  <i class="fas fa-download"></i>
                </span>
                <span>Download</span>
              </a>
              </span>
              <span> configuration file</span></p>
            <p>Open VPN client and click plus icon to add a new tunnel. Choose import tunnel from file, then locate the configuration file you just downloaded to import.</p>
            <p>If your newly added tunnel status is inactive, click activate button to activate it.</p>
          </div>
        </div>
        <div class="is-divider-vertical" data-content="OR"></div>
        <div class="column has-text-centered">
          <div class="content">
            <h3>For Mobile (Android/iOS/ipad)</h3>
            <p>Open VPN client and scan this QRCode:</p>
            <img alt="Device QR code" src="data:image/png;base64,{{ $.qrcode }}" />
          </div>
        </div>
      </div>
      <div class="is-divider" data-content="Step three"></div>
      <div id="step-three" class="columns is-centered">
        <div class="column has-text-centered">
          <div class="content">
            <h3>Check you tunnel status</h3>
            <p>The index page list your created devices and its status, click to return to <a href="/">index page</a>.</p>
          </div>
        </div>
      </div>
    </section>
    </div> <!--container-->
  </body>
</html>
