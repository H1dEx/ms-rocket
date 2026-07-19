package payment

import "regexp"

func (s *ServiceSuite) TestPayOrderSuccess() {
	uuid, err := s.service.PayOrder(s.ctx)

	s.Require().NoError(err)
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	s.Regexp(uuidRegex, uuid)
}
