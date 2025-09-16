import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DodajStudentaUSobu } from './dodaj-studenta-usobu';

describe('DodajStudentaUSobu', () => {
  let component: DodajStudentaUSobu;
  let fixture: ComponentFixture<DodajStudentaUSobu>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DodajStudentaUSobu]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DodajStudentaUSobu);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
