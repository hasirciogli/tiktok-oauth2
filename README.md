# TikTok OAuth2 v2 Server

Bu proje, TikTok OAuth2 v2 API'sini kullanarak kullanıcı kimlik doğrulaması yapan Go tabanlı REST API sunucusudur.

## Özellikler

- ✅ TikTok OAuth2 v2 desteği
- ✅ Authorization code flow
- ✅ Token refresh mekanizması
- ✅ CSRF koruması (state parameter)
- ✅ CORS desteği
- ✅ JSON API responses
- ✅ Error handling
- ✅ Environment variable configuration

## Kurulum

1. **Bağımlılıkları yükle:**
```bash
go mod tidy
```

2. **Environment variables ayarla:**
```bash
# .env dosyası oluştur
cp .env.example .env

# .env dosyasını düzenle
TIKTOK_CLIENT_KEY=your_client_key_here
TIKTOK_CLIENT_SECRET=your_client_secret_here
TIKTOK_REDIRECT_URI=https://yourdomain.com/callback
SERVER_PORT=8080
```

3. **Sunucuyu başlat:**
```bash
go run main.go
```

## API Endpoints

### 1. Health Check
```
GET /health
```
Sunucu durumunu kontrol eder.

### 2. OAuth Authorization
```
GET /auth
```
Kullanıcıyı TikTok OAuth sayfasına yönlendirir.

### 3. OAuth Callback
```
GET /callback?code=AUTH_CODE&state=STATE
```
TikTok'dan gelen authorization code'u access token ile değiştirir.

### 4. Token Refresh
```
POST /refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```
Access token'ı yeniler.

### 5. User Info
```
GET /user
Authorization: Bearer YOUR_ACCESS_TOKEN
```
Kullanıcı bilgilerini getirir (access token gerekli).

## Kullanım

1. **OAuth flow başlat:**
   ```
   GET http://localhost:8080/auth
   ```

2. **TikTok'da izin ver ve callback'e yönlendiril**

3. **Token ve kullanıcı bilgilerini al:**
   ```json
   {
     "success": true,
     "message": "Authentication successful",
     "data": {
       "token": {
         "access_token": "act.xxx...",
         "expires_in": 86400,
         "open_id": "user_open_id",
         "refresh_token": "rft.xxx...",
         "refresh_expires_in": 31536000
       },
       "user_info": {
         "open_id": "user_open_id",
         "union_id": "user_union_id",
         "avatar_url": "https://...",
         "avatar_url_100": "https://...",
         "avatar_large_url": "https://...",
         "display_name": "User Display Name",
         "bio_description": "User bio...",
         "profile_deep_link": "https://...",
         "is_verified": false,
         "username": "username",
         "follower_count": 1000,
         "following_count": 500,
         "likes_count": 10000,
         "video_count": 50
       }
     }
   }
   ```

4. **Ayrıca kullanıcı bilgilerini tekrar al:**
   ```
   GET /user
   Authorization: Bearer YOUR_ACCESS_TOKEN
   ```

## TikTok Developer Setup

1. [TikTok for Developers](https://developers.tiktok.com/) hesabı oluştur
2. Yeni bir uygulama oluştur
3. OAuth2 v2 izinlerini etkinleştir
4. Redirect URI'yi ayarla (HTTPS gerekli)
5. Client Key ve Client Secret'ı al
6. **Gerekli Scope'lar:**
   - `user.info.basic` - Temel kullanıcı bilgileri
   - `user.info.profile` - Profil bilgileri
   - `user.info.stats` - İstatistik bilgileri
   - `video.list` - Video listesi (opsiyonel)

## Güvenlik Notları

- Production'da HTTPS kullanın
- State parameter'ı doğru şekilde validate edin
- Client secret'ı güvenli tutun
- Token'ları güvenli şekilde saklayın

## Geliştirme

```bash
# Bağımlılıkları güncelle
go mod tidy

# Test et
go run main.go

# Build et
go build -o tiktok-oauth2 main.go
```

## Railway Deployment

### 1. Railway'a Deploy Et

```bash
# Railway CLI ile
railway login
railway init
railway up

# Veya GitHub ile bağla
# Railway dashboard'dan GitHub repo'yu seç
```

### 2. Environment Variables Ayarla

Railway dashboard'da şu değişkenleri ekle:

```
TIKTOK_CLIENT_KEY=your_client_key
TIKTOK_CLIENT_SECRET=your_client_secret
TIKTOK_REDIRECT_URI=https://your-app.railway.app/callback
SERVER_PORT=8080
```

### 3. Build & Start Commands

Railway otomatik olarak algılar:
- **Build:** `go build -o tiktok-oauth2 main.go`
- **Start:** `./tiktok-oauth2`
- **Health Check:** `GET /health`

### 4. Docker ile Deploy

```bash
# Docker build
docker build -t tiktok-oauth2 .

# Docker run
docker run -p 8080:8080 \
  -e TIKTOK_CLIENT_KEY=your_key \
  -e TIKTOK_CLIENT_SECRET=your_secret \
  -e TIKTOK_REDIRECT_URI=https://yourdomain.com/callback \
  tiktok-oauth2
```

## Lisans

MIT License
