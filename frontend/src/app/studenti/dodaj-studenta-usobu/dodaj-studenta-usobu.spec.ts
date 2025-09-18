import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DodajStudentaUSobuComponent } from './dodaj-studenta-usobu';

describe('DodajStudentaUSobu', () => {
  let component: DodajStudentaUSobuComponent;
  let fixture: ComponentFixture<DodajStudentaUSobuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DodajStudentaUSobuComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DodajStudentaUSobuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
