---

# Authentication API Endpoints

## Base URL:
`http://localhost:8085/api/auth`

### 1. **Registration**
**POST** `http://localhost:8085/api/auth/register`

#### Request:
```json
{
  "fullname": "Swan Htet Aung",
  "email": "swanhtetaungp@gmail.com",
  "country": "USA",
  "biometric_hash": "123456dfafsdiiiads789abcdef",
  "password": "securepassword123",
  "timestamp": 1679987456
}
```

#### Response:
```json
{
  "message": "Registration succeed"
}
```

---

### 2. **Verify Email**
**POST** `http://localhost:8085/api/auth/verify`

#### Request:
```json
{
  "email": "swanhtetaungp@gmail.com",
  "code": "257865"
}
```

#### Response:
```json
{
  "message": "Verification succeeded"
}
```

---

### 3. **Login**
**POST** `http://localhost:8085/api/auth/login`

#### Request:
```json
{
  "email": "swanhtetaungp@gmail.com",
  "password": "securepassword123"
}
```

#### Response:
```json
{
  "message": "Login successful",
  "token": "your-jwt-token"
}
```

---

### 4. **Wallet**
**POST** `http://localhost:8085/api/auth/wallet`

#### Headers:
- Authorization: Bearer `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

#### Request:
```json
{
  "userid": "swanhtetaungp@gmail.com",
  "phrase": "fdajfl;adsjfl;adsd"
}
```

#### Response:
```json
{
  "message": "Wallet data successfully retrieved"
}
```

---

### 5. **Get User Information (Me)**
**POST** `http://localhost:8085/api/auth/me?email=swanhtetaungp@gmail.com`

#### Headers:
- Authorization: Bearer `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

#### Request:
No body required.

#### Response:
```json
{
  "email": "swanhtetaungp@gmail.com",
  "fullname": "Swan Htet Aung",
  "country": "USA"
}
```

---

### 6. **Refresh Token**
**GET** `http://localhost:8085/api/auth/refresh`

#### Headers:
- Cookie: refresh_token=`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

#### Response:
```json
{
  "message": "Token refreshed successfully",
  "new_token": "new-jwt-token"
}
```

--- 
