package uploader

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (c *Client) browserLogin(page *rod.Page) error {
	p := page.Timeout(c.actionTimeout)

	if err := p.Navigate(c.serverURL + "/public"); err != nil {
		return fmt.Errorf("navigate to public: %w", err)
	}

	loginLink, err := p.ElementR("a", "Login")
	if err != nil {
		return fmt.Errorf("find login link: %w", err)
	}
	if err := loginLink.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click login link: %w", err)
	}

	if _, err := p.Element("input[placeholder=\"Username\"]"); err != nil {
		return fmt.Errorf("wait for login form: %w", err)
	}

	if err := fillByPlaceholder(p, "Username", c.username); err != nil {
		return err
	}
	if err := fillByPlaceholder(p, "Password", c.password); err != nil {
		return err
	}

	waitNav := page.WaitNavigation(proto.PageLifecycleEventNameNetworkIdle)
	if err := clickSubmit(p); err != nil {
		return fmt.Errorf("login submit: %w", err)
	}
	waitNav()

	return nil
}

func (c *Client) browserUpload(page *rod.Page, uploadURL, inputFile, outputFile string) error {
	p := page.Timeout(c.actionTimeout)

	if err := p.Navigate(uploadURL); err != nil {
		return fmt.Errorf("navigate to upload page: %w", err)
	}

	if _, err := p.ElementR("label", "Input testdata"); err != nil {
		return fmt.Errorf("wait for upload form: %w", err)
	}

	if err := setFileByLabel(p, "Input testdata", inputFile); err != nil {
		return err
	}
	if err := setFileByLabel(p, "Output testdata", outputFile); err != nil {
		return err
	}

	waitNav := page.WaitNavigation(proto.PageLifecycleEventNameNetworkIdle)
	if err := clickInputSubmit(p, "Submit all changes"); err != nil {
		return fmt.Errorf("upload submit: %w", err)
	}
	waitNav()

	return nil
}

func clickSubmit(p *rod.Page) error {
	if el, err := p.Element(`button[type="submit"]`); err == nil {
		return el.Click(proto.InputMouseButtonLeft, 1)
	}
	el, err := p.Element(`input[type="submit"]`)
	if err != nil {
		return fmt.Errorf("no submit element found: %w", err)
	}
	return el.Click(proto.InputMouseButtonLeft, 1)
}

func clickInputSubmit(p *rod.Page, value string) error {
	el, err := p.Element(fmt.Sprintf(`input[type="submit"][value="%s"]`, value))
	if err != nil {
		el, err = p.Element(`input[type="submit"]`)
		if err != nil {
			return fmt.Errorf("no submit input found: %w", err)
		}
	}
	return el.Click(proto.InputMouseButtonLeft, 1)
}

func fillByPlaceholder(p *rod.Page, placeholder, value string) error {
	el, err := p.Element(fmt.Sprintf(`input[placeholder="%s"]`, placeholder))
	if err != nil {
		return fmt.Errorf("find placeholder=%q: %w", placeholder, err)
	}
	if err := el.Input(value); err != nil {
		return fmt.Errorf("fill placeholder=%q: %w", placeholder, err)
	}
	return nil
}

func setFileByLabel(p *rod.Page, labelText, filePath string) error {
	label, err := p.ElementR("label", labelText)
	if err != nil {
		return fmt.Errorf("find label %q: %w", labelText, err)
	}
	forAttr, err := label.Attribute("for")
	if err != nil || forAttr == nil {
		return fmt.Errorf("label %q missing 'for' attribute", labelText)
	}
	input, err := p.Element(fmt.Sprintf(`input#%s`, *forAttr))
	if err != nil {
		return fmt.Errorf("find input for label %q: %w", labelText, err)
	}
	if err := input.SetFiles([]string{filePath}); err != nil {
		return fmt.Errorf("set file for label %q: %w", labelText, err)
	}
	return nil
}
