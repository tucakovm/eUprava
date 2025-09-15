import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface CanteenDto {
  id: string;
  name: string;
  address: string;
  open_at: string;
  close_at: string;
}

@Injectable({
  providedIn: 'root'
})
export class CanteenService {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8001/api/canteens/';

  getAll(): Observable<CanteenDto[]> {
    return this.http.get<CanteenDto[]>(this.baseUrl);
  }

  getOne(id: string): Observable<CanteenDto> {
    return this.http.get<CanteenDto>(`${this.baseUrl}${id}`);
  }

  delete(id: string) {
    return this.http.delete(`${this.baseUrl}${id}`);
  }

  create(canteen: CanteenDto) {
    return this.http.post<CanteenDto>(`${this.baseUrl}`, canteen);
  }


}
