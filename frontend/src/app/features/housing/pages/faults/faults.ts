import { Component } from '@angular/core';
import { CommonModule, JsonPipe, NgIf } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators, FormGroup } from '@angular/forms';
import { HousingService } from '../../../../services/housing.service';
import { Kvar, StatusKvara } from '../../../../model/housing';

@Component({
  selector: 'app-faults',
  standalone: true,
   imports: [CommonModule, ReactiveFormsModule, NgIf, JsonPipe],
  templateUrl: './faults.html'
})
export class Faults {
  created?: Kvar;
  statusRes?: { status: string };
  error?: string;

  reportForm!: FormGroup;
  statusForm!: FormGroup;

  constructor(private fb: FormBuilder, private api: HousingService) {
    this.reportForm = this.fb.group({
      sobaId: ['', Validators.required],
      prijavioId: ['', Validators.required],
      opis: ['', Validators.required]
    });

    this.statusForm = this.fb.group({
      kvarId: ['', Validators.required],
      status: ['u_toku', Validators.required]
    });
  }

  report() {
    if (this.reportForm.invalid) return;
    const { sobaId, prijavioId, opis } = this.reportForm.value;
    this.api.reportFault(sobaId!, prijavioId!, opis!).subscribe({
      next: k => this.created = k,
      error: e => this.error = e?.error || 'Error'
    });
  }

  changeStatus() {
    if (this.statusForm.invalid) return;
    const { kvarId, status } = this.statusForm.value;
    this.api.changeFaultStatus(kvarId!, status as StatusKvara).subscribe({
      next: res => this.statusRes = res,
      error: e => this.error = e?.error || 'Error'
    });
  }
}
