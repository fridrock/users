package registration

import "net/http"

type RegistrationHandler interface {
	HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error)
}

type RegistrationHandlerImpl struct {
	storage RegistrationStorage
	parser  RegistrationParser
}

//TODO response with token
//TODO make validation
//TODO tests
func (rs *RegistrationHandlerImpl) HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error) {
	userDto, err := rs.parser.GetDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	id, err := rs.storage.SaveUser(userDto)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte(id.String()))
	return http.StatusCreated, nil
}
func NewRegistrationHandler(storage RegistrationStorage) RegistrationHandler {
	return &RegistrationHandlerImpl{
		storage: storage,
		parser:  &RegistrationParserImpl{},
	}
}
