import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OstaviRecenzijuComponent } from './ostavi-recenziju.component';

describe('OstaviRecenzijuComponent', () => {
  let component: OstaviRecenzijuComponent;
  let fixture: ComponentFixture<OstaviRecenzijuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OstaviRecenzijuComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OstaviRecenzijuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
