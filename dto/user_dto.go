/**
@filename:dto
@author: silicon-he
@data:2021/6/19
@software:Goland
@note:
**/

package dto

import "silicon.com/ginessential/model"

type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User)UserDto  {
	return UserDto{
		Name: user.Name,
		Telephone: user.Telephone,
	}
}
