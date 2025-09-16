import {inject, Injectable} from '@angular/core';
import {CanteenDto} from './canteen.service';
import {HttpClient} from '@angular/common/http';
import {Menu} from '../model/menus';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MenuService {

  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8001/api/menus/';

  create(menu: Menu) {
    return this.http.post<Menu>(`${this.baseUrl}`, menu);
  }

  getAll(canteenId: string): Observable<Menu[]> {
    return this.http.get<Menu[]>(`${this.baseUrl}${canteenId}`);
  }

  delete(id: string) {
    return this.http.delete(`${this.baseUrl}${id}`);
  }
}
