import msgpack
import requests

'''data = {
    "Login": "tux",
    "Password": "12345678",
    "FirstName": "Tux",
    "LastName": "Torvalds",
    "Patronymic": "-",
}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/user_create', data=query)

print(response.status_code)
print(response.content)'''
data = {
    "Token":       "CCmS1zwIRIukZwii31xOwrkrA2cz+fWA",
	"Name":        "Third",
	"Description": "",
	"District":    "",
	"Region":      "",
	"Address":    "",
	"Type":        "",
	"State":       "",
	"AreaFrom":        10,
    "AreaTo":      10003,
	"Owner":       "",
	"Actual_user": "",
	"Gid":         1,
    "Limit":       2,
	"SortAsc":     False

}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/find_object', data=query)

print(response.status_code)
print(msgpack.unpackb(response.content))
'''
data = {
    "Token": "CCmS1zwIRIukZwii31xOwrkrA2cz+fWA",
    "Name": "Main"
}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/group_create', data=query)

print(response.status_code)
print(msgpack.unpackb(response.content))'''
