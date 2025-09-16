import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import {
  Student, Soba, RecenzijaSobe, Kvar, StatusKvara, StudentskaKartica
} from '../model/housing';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class HousingService {
  private base = 'http://localhost:8003/api/housing';

  constructor(private http: HttpClient) {}

  // Students
  createStudent(ime: string, prezime: string): Observable<Student> {
    return this.http.post<Student>(`${this.base}/students`, { ime, prezime });
  }

  releaseStudentRoom(studentId: string): Observable<{ status: string }> {
    return this.http.post<{ status: string }>(`${this.base}/students/release`, { studentId });
  }

  // Studentska kartica
  createStudentCardIfMissing(studentId: string): Observable<StudentskaKartica> {
    return this.http.post<StudentskaKartica>(`${this.base}/students/cards`, { studentId });
  }

  getStudentCard(studentId: string): Observable<StudentskaKartica> {
    const params = new HttpParams().set('studentId', studentId);
    return this.http.get<StudentskaKartica>(`${this.base}/students/cards`, { params });
  }

  updateStudentCardBalance(studentId: string, delta: number): Observable<StudentskaKartica> {
    return this.http.post<StudentskaKartica>(`${this.base}/students/cards/balance`, { studentId, delta });
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

  assignStudentToRoom(domId: string, broj: string, ime: string, prezime: string): Observable<Student> {
    return this.http.post<Student>(`${this.base}/rooms/assign`, { domId, broj, ime, prezime });
  }

  listFreeRooms(domId: string): Observable<Soba[]> {
    const params = new HttpParams().set('domId', domId);
    return this.http.get<Soba[]>(`${this.base}/rooms/free`, { params });
  }

  // Reviews
  addRoomReview(sobaId: string, autorId: string, ocena: number, komentar?: string | null):
    Observable<RecenzijaSobe> {
    return this.http.post<RecenzijaSobe>(`${this.base}/rooms/reviews`,
      { sobaId, autorId, ocena, komentar: komentar ?? null });
  }

  // Faults
  reportFault(sobaId: string, prijavioId: string, opis: string): Observable<Kvar> {
    return this.http.post<Kvar>(`${this.base}/rooms/faults`, { sobaId, prijavioId, opis });
  }

  changeFaultStatus(kvarId: string, status: StatusKvara): Observable<{ status: string }> {
    return this.http.post<{ status: string }>(`${this.base}/faults/status`, { kvarId, status });
  }
}
