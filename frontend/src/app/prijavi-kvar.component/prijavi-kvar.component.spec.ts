import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PrijaviKvarComponent } from './prijavi-kvar.component';

describe('PrijaviKvarComponent', () => {
  let component: PrijaviKvarComponent;
  let fixture: ComponentFixture<PrijaviKvarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PrijaviKvarComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PrijaviKvarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
