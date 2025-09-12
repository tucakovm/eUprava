# === Build (deps) stage ===
FROM node:20-alpine AS deps
WORKDIR /app

# Kopiraj samo package fajlove i instaliraj zavisnosti
COPY package*.json ./
RUN npm ci

# Kopiraj ostatak koda
COPY . .

# === Runtime stage ===
FROM node:20-alpine
WORKDIR /app

# Prenesi node_modules i kod iz deps stage-a
COPY --from=deps /app /app

# Dodaj lokalne bin-ove u PATH (tu je i "ng")
ENV PATH=/app/node_modules/.bin:$PATH

EXPOSE 4200

# Pokreni lokalni Angular CLI (bez globalnog)
CMD ["ng", "serve", "--host", "0.0.0.0", "--port", "4200", "--poll", "2000"]
