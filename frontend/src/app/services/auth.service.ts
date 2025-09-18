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
            localStorage.setItem('user', res.user.id); // direktan string
            localStorage.setItem('role', res.user.role); // direktan string
          }
          return res;
        }),
        catchError(this.handle)
      );
  }

  get userRole(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('role'); // vrati direktno string
  }

  get userId(): string | null {
    if (typeof window === 'undefined') return null;
    const id = localStorage.getItem('user');
    return id && id !== 'null' && id !== 'undefined' ? id : null;
  }


  get token(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  get isAuthenticated(): boolean {
    return !!this.token;
  }

  logout() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      localStorage.removeItem('role');
    }
  }

  private handle(err: HttpErrorResponse) {
    const msg = err.error?.error || err.message || 'Došlo je do greške';
    return throwError(() => new Error(msg));
  }
}
