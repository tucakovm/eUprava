import { Component } from '@angular/core';
import { CommonModule, JsonPipe, NgIf, NgFor } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormGroup } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { HousingService } from '../../../../services/housing.service';
import { Soba, Student } from '../../../../model/housing';
import { Router } from '@angular/router';

@Component({
  selector: 'app-rooms',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, NgIf, NgFor, JsonPipe],
  templateUrl: './rooms.html'
})
export class Rooms {
  assigned?: Student;
  room?: Soba;
  freeRooms: Soba[] = [];
  error?: string;

  assignForm!: FormGroup;
  freeForm!: FormGroup;
  getForm!: FormGroup;

  constructor(private fb: FormBuilder, private api: HousingService, private router: Router) {
    this.assignForm = this.fb.group({
      domId: ['', Validators.required],
      broj: ['', Validators.required],
      ime: ['', Validators.required],
      prezime: ['', Validators.required]
    });

    this.freeForm = this.fb.group({
      domId: ['', Validators.required]
    });

    this.getForm = this.fb.group({
      id: ['', Validators.required]
    });
  }

  submitAssign() {
    this.error = undefined;
    if (this.assignForm.invalid) return;
    const v = this.assignForm.value;
    this.api.assignStudentToRoom(v.domId!, v.broj!, v.ime!, v.prezime!).subscribe({
      next: s => this.assigned = s,
      error: e => this.error = e?.error || 'Error'
    });
  }

  submitFree() {
    this.error = undefined;
    if (this.freeForm.invalid) return;
    const { domId } = this.freeForm.value;
    this.api.listFreeRooms(domId!).subscribe({
      next: rooms => this.freeRooms = rooms,
      error: e => this.error = e?.error || 'Error'
    });
  }

  submitGet() {
    this.error = undefined;
    if (this.getForm.invalid) return;
    const { id } = this.getForm.value;
    this.api.getRoom(id!).subscribe({
      next: r => this.room = r,
      error: () => this.error = 'Not found'
    });
  }

  detail(id: string) {
    this.router.navigate(['/housing/rooms', id]);
  }
}
