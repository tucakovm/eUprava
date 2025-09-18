import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SlobodneSobe } from './slobodne-sobe';

describe('SlobodneSobe', () => {
  let component: SlobodneSobe;
  let fixture: ComponentFixture<SlobodneSobe>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SlobodneSobe]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SlobodneSobe);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
