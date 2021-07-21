
func NewCredentials(certFile, serverNameOverride, keyFile string) (credentials.TransportCredentials, error) {
	cp := x509.NewCertPool()
	if certFile != "" {
		b, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
	}

	w, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{
		ServerName: serverNameOverride,
		//RootCAs:      cp,
		KeyLogWriter: w,
	}), nil
}

