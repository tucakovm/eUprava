import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { RegisterRequest, LoginRequest, LoginResponse, User } from '../model/auth';
import { catchError, map, throwError } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8002';

  register(body: RegisterRequest) {
    return this.http.post<User>(`${this.baseUrl}/api/register`, body)
      .pipe(catchError(this.handle));
  }

  login(body: LoginRequest) {
    return this.http.post<LoginResponse>(`${this.baseUrl}/api/login`, body)
      .pipe(
        map(res => {
          if (typeof window !== 'undefined') {
            localStorage.setItem('token', res.token);
            localStorage.setItem('user', JSON.stringify(res.user.username));
          }
          return res;
        }),
        catchError(this.handle)
      );
  }

  logout() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  }

  get token(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  get isAuthenticated(): boolean {
    return !!this.token;
  }

  private handle(err: HttpErrorResponse) {
    const msg = err.error?.error || err.message || 'Došlo je do greške';
    return throwError(() => new Error(msg));
  }
}
