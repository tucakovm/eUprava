import { Component } from '@angular/core';
import { CommonModule, JsonPipe, NgIf } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormGroup } from '@angular/forms';
import { HousingService } from '../../../../services/housing.service';
import { Student } from '../../../../model/housing';

@Component({
  selector: 'app-students',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, NgIf, JsonPipe],
  templateUrl: './students.html'
})
export class Students {
  created?: Student;
  releasedStatus?: string;
  error?: string;

  createForm!: FormGroup;
  releaseForm!: FormGroup;

  constructor(private fb: FormBuilder, private api: HousingService) {
    this.createForm = this.fb.group({
      ime: ['', Validators.required],
      prezime: ['', Validators.required],
    });

    this.releaseForm = this.fb.group({
      studentId: ['', Validators.required],
    });
  }

  submitCreate() {
    this.error = undefined;
    if (this.createForm.invalid) return;
    const { ime, prezime } = this.createForm.value;
    this.api.createStudent(ime!, prezime!).subscribe({
      next: s => this.created = s,
      error: e => this.error = e?.error || 'Error'
    });
  }

  submitRelease() {
    this.error = undefined;
    if (this.releaseForm.invalid) return;
    const { studentId } = this.releaseForm.value;
    this.api.releaseStudentRoom(studentId!).subscribe({
      next: res => this.releasedStatus = res.status,
      error: e => this.error = e?.error || 'Error'
    });
  }
}
