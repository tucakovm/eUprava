import { Injectable } from '@angular/core';
import { HttpClient, HttpParams , HttpHeaders} from '@angular/common/http';
import {Dom, Student, Soba, RecenzijaSobe, Kvar, StatusKvara, StudentskaKartica
} from '../model/housing';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class HousingService {
  private base = 'http://localhost:8003/api/housing';

  constructor(private http: HttpClient , private auth:AuthService) {}

    getAllDoms() {
    return this.http.get<Dom[]>(`${this.base}/doms`);
  }

  getDomById(id: string) {
    const params = new HttpParams().set('id', id);
    return this.http.get<Dom>(`${this.base}/dom`, { params });
  }

  // Students
  createStudent(ime: string, prezime: string): Observable<Student> {
    return this.http.post<Student>(`${this.base}/students`, { ime, prezime });
  }

  releaseStudentRoom(studentId: string): Observable<{ status: string }> {
    return this.http.post<{ status: string }>(`${this.base}/students/release`, { studentId });
  }

  // Studentska kartica
  createStudentCardIfMissing(studentUsername: string): Observable<StudentskaKartica> {
    return this.http.post<StudentskaKartica>(`${this.base}/students/cards`, { studentUsername });
  }

  getStudentCard(studentUsername: string): Observable<StudentskaKartica> {
    const params = new HttpParams().set('studentId', studentUsername);
    return this.http.get<StudentskaKartica>(`${this.base}/students/cards`, { params });
  }

  updateStudentCardBalance(studentUsername: string, delta: number): Observable<StudentskaKartica> {
    return this.http.post<StudentskaKartica>(`${this.base}/students/cards/balance`, { studentUsername, delta });
  }

  // Rooms
  getRoom(id: string): Observable<Soba> {
    const params = new HttpParams().set('id', id);
    return this.http.get<Soba>(`${this.base}/rooms`, { params });
  }

  getRoomDetail(id: string): Observable<Soba> {
    const params = new HttpParams().set('id', id);
    return this.http.get<Soba>(`${this.base}/rooms/detail`, { params });
  }

  assignStudentToRoom(domId: string, broj: string, username: string) {
    return this.http.post<Student>(`${this.base}/rooms/assign`, { domId, broj, username });
  }

  listFreeRooms(domId: string): Observable<Soba[]> {
    const params = new HttpParams().set('domId', domId);
    return this.http.get<Soba[]>(`${this.base}/rooms/free`, { params });
  }

  addRoomReview(
  sobaId: string,
  autorUsername: string,
  ocena: number,
  komentar?: string | null
): Observable<RecenzijaSobe> {
  const body = {
    sobaId,
    autorUsername,
    ocena,
    komentar: (komentar ?? '').trim() || null,
  };

  const headers = {
    'Content-Type': 'application/json',
    ...(this.auth.token ? { Authorization: `Bearer ${this.auth.token}` } : {}),
  };

  return this.http.post<RecenzijaSobe>(
    `${this.base}/rooms/reviews`,
    body,
    { headers }
  );
}



  // Faults
  reportFault(sobaId: string, prijavioUsername: string, opis: string): Observable<Kvar> {
    return this.http.post<Kvar>(`${this.base}/rooms/faults`, { sobaId, prijavioUsername, opis });
  }

  changeFaultStatus(kvarId: string, status: StatusKvara): Observable<{ status: string }> {
    return this.http.post<{ status: string }>(`${this.base}/faults/status`, { kvarId, status });
  }
}
