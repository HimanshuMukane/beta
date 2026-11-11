import {
  Card,
  CardHeader,
  CardContent,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { toast } from '@repo/ui/components/ui/sonner'
import {
  IconMessageHeart,
  IconStarFilled,
  IconMoodSmile,
  IconThumbUp,
  IconFileSpreadsheet,
} from "@repo/ui/icons/tabler-icons"

export function FeedbackReportPage() {
  const reviews = [
    {
      id: 'rev-1',
      client: 'Kavya & Siddharth',
      event: 'Reception & Wedding Ceremony',
      rating: 5,
      date: '2026-11-05',
      photographer: 'Rohan Mehta & Team',
      npsScore: 10,
      comment:
        '“The team was exceptionally punctual and caught all the subtle, candid family emotions without intruding. Our teaser portraits blew everyone away on Instagram!”',
      tags: ['Candid Mastery', 'Super Punctual', 'Great Teasers'],
    },
    {
      id: 'rev-2',
      client: 'Deepika & Rohit Rao',
      event: 'Maternity Studio Portrait',
      rating: 5,
      date: '2026-11-04',
      photographer: 'Priya Sharma',
      npsScore: 10,
      comment:
        '“Studio ambiance was warm and comforting. Priya guided my poses very patiently. Photo delivery happened 2 days ahead of the promised date.”',
      tags: ['Comfortable Atmosphere', 'Early Delivery'],
    },
    {
      id: 'rev-3',
      client: 'Arjun Verma',
      event: 'Sangeet Night Live Coverage',
      rating: 5,
      date: '2026-11-02',
      photographer: 'Aditya Roy',
      npsScore: 9,
      comment:
        '“Lighting on the dance floor was extremely challenging, but the low-light prime lens captures were sharp and colorful. Truly world-class work.”',
      tags: ['Low Light Expertise', 'Vibrant Colors'],
    },
    {
      id: 'rev-4',
      client: 'Rohan Kapoor',
      event: 'Goa Destination Pre-Wedding',
      rating: 4,
      date: '2026-10-20',
      photographer: 'Vikram Seth (Drone) & Rohan',
      npsScore: 8,
      comment:
        '“The drone sunset aerials are breathtaking. Album layout design was stunning, though we had a 1-day delay during courier delivery.”',
      tags: ['Drone Perfection', 'Luxury Photobook'],
    },
  ]

  const ratingStats = {
    average: 4.9,
    totalReviews: 86,
    npsScore: 88,
    csatPercentage: 97.4,
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-rose-600 text-white shadow-2xs">
              <IconMessageHeart className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Client Feedback & NPS Satisfaction Report
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Verified client testimonials, Net Promoter Score (NPS), and photographer ratings.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Reviews exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-rose-600" />
            <span>Export Reviews</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Overall Client Rating
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground flex items-center gap-1.5">
              <span>{ratingStats.average}</span>
              <div className="flex text-amber-500">
                <IconStarFilled className="h-5 w-5" />
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">Based on {ratingStats.totalReviews} verified shoot surveys</p>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Net Promoter Score (NPS)
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              +{ratingStats.npsScore}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconThumbUp className="h-3.5 w-3.5" />
              <span>World-class client loyalty</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              CSAT Satisfaction
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-primary">
              {ratingStats.csatPercentage}%
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconMoodSmile className="h-3.5 w-3.5 text-primary" />
              <span>Rated 4 or 5 stars out of 5</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Client Testimonials
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-purple-600">
              100%
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-xs text-muted-foreground">Public permission granted for portfolio</p>
          </CardContent>
        </Card>
      </div>

      {/* Review Cards Feed */}
      <div className="space-y-4">
        <h3 className="text-sm font-bold text-foreground tracking-tight">Recent Client Reviews & Reflections</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {reviews.map((rev) => (
            <Card key={rev.id} className="border bg-card shadow-2xs p-4 flex flex-col justify-between gap-3">
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <h4 className="font-bold text-sm text-foreground">{rev.client}</h4>
                    <p className="text-[11px] text-muted-foreground">{rev.event} • {rev.photographer}</p>
                  </div>
                  <div className="flex items-center gap-1 text-amber-500">
                    {Array.from({ length: rev.rating }).map((_, i) => (
                      <IconStarFilled key={i} className="h-3.5 w-3.5" />
                    ))}
                  </div>
                </div>

                <p className="text-xs text-foreground/90 italic leading-relaxed">
                  {rev.comment}
                </p>
              </div>

              <div className="flex flex-wrap items-center gap-1.5 pt-2 border-t text-[10px]">
                {rev.tags.map((t) => (
                  <Badge key={t} variant="secondary" className="text-[10px] py-0.5">
                    {t}
                  </Badge>
                ))}
                <span className="ml-auto text-muted-foreground">{rev.date}</span>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}
